package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/internal/ArgsParser"
	"go-speedtest-bot/internal/config"
	"go-speedtest-bot/internal/database"
	"go-speedtest-bot/internal/speedtest"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadConf() *Conf {
	cfg := config.GetBotFile()
	s := cfg.Section("bot")
	return &Conf{token: s.Key("token").String()}
}

// NewBot return a bot instance
func NewBot() *B {
	bot, err := tgbotapi.NewBotAPI(loadConf().token)
	if err != nil {
		log.Println("[NewBotError]Token", err)
		os.Exit(-1)
	}
	return bot
}

// SendT send text message
func SendT(bot *B, cid int64, text string) {
	msg := tgbotapi.NewMessage(cid, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("SendError", err)
	}
}

// SendP send parsed text message
func SendP(bot *B, cid int64, text string, format string) {
	msg := tgbotapi.NewMessage(cid, text)
	msg.ParseMode = format
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("SendError", err)
	}
}

// Launch is the robot's main methods. Calling this method will
// enter a loop and continuously listen for messages.
// If debug, bot will enable debug mode.
// If logInfo, program will print all the message.
// If clean, program will clean all the out of date message.
func Launch(debug bool, logInfo bool, clean bool) {
	bot := NewBot()
	bot.Debug = debug
	log.Println("Authorized on account", bot.Self.UserName)

	err := LoadAdmin()
	if err != nil {
		log.Println("Fail to load admin list.", err)
		os.Exit(-1)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	if clean {
		log.Println("Cleaning mode on")
		updates.Clear()
		log.Println("message all clear.")
		os.Exit(0)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if logInfo {
			log.Printf("[%s]%s", update.Message.From.UserName, update.Message.Text)
		}
		if Auth(int64(update.Message.From.ID)) {
			CMDHandler(bot, update.Message)
		}
	}
}

// CMDHandler handle all the command
func CMDHandler(bot *B, msg *M) {
	if msg.IsCommand() {
		if cmd, ok := Commands[msg.Command()]; ok {
			cmd(bot, msg)
		}
	}
}

// cmd /start
func cmdStart(b *B, m *M) {
	text := "Here is a bot who can help you manage all your proxy."
	SendT(b, m.Chat.ID, text)
}

// cmd /ping
func cmdPing(b *B, m *M) {
	connected := speedtest.Ping(speedtest.GetHost())
	var text string
	if connected {
		text = "Connect to backend successfully"
	} else {
		text = "Unable to connect to the backend, please check out the latest logs."
	}
	SendT(b, m.Chat.ID, text)
}

// cmd /status
func cmdStatus(b *B, m *M) {
	result, err := speedtest.GetStatus(speedtest.GetHost())
	if err != nil {
		SendT(b, m.Chat.ID, fmt.Sprint(err))
		return
	}
	SendT(b, m.Chat.ID, "Status: "+result.State)
}

// cmd /read_sub
func cmdReadSub(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 2 {
		SendT(b, m.Chat.ID, "Use case(Only single link is supported):\n/read_sub https://example.com")
		return
	}
	url := strings.Fields(m.Text)[1]
	subResps, err := speedtest.ReadSubscriptions(speedtest.GetHost(), url)
	if err != nil {
		SendT(b, m.Chat.ID, fmt.Sprint(err))
		return
	}
	var text string
	for _, subResp := range subResps {
		if subResp.Type != "" && subResp.Config != nil {
			text += subResp.Type + " " + subResp.Config.Remarks + "\n"
		}
	}
	SendT(b, m.Chat.ID, text)
}

func formatResult(r *speedtest.Result) string {
	if len(r.Result) == 0 {
		return ""
	}
	var text string
	text = "Recent result(ls=loss, lp=local ping, gp=google ping):\n"
	text += "Status: " + r.Status + "\n"
	if len(r.Current.Remarks) != 0 {
		text += "Nodes being tested: " + r.Current.Remarks + "\n"
	}
	text += "\n"
	for _, res := range r.Result {
		text += fmt.Sprintf(
			"%s: | ls: %.2f%% | lp: %.2f ms | gp: %.2f ms\n", res.Remarks, res.Loss*100, res.Ping*1000, res.GPing*1000)
	}
	return text
}

// cmd /result
func cmdResult(b *B, m *M) {
	result, err := speedtest.GetResult(speedtest.GetHost())
	if err != nil {
		SendT(b, m.Chat.ID, fmt.Sprint(err))
		return
	}
	if fresult := formatResult(result); len(fresult) != 0 {
		SendT(b, m.Chat.ID, fresult)
		return
	}
	SendT(b, m.Chat.ID, "No result yet")
}

func startTestWithURL(b *B, m *M, url string, method string, mode string, include []string, exclude []string) {
	result, err := speedtest.GetStatus(speedtest.GetHost())
	if err != nil {
		SendT(b, m.Chat.ID, err.Error())
		return
	}
	if result == nil {
		SendT(b, m.Chat.ID, "Unable to fetch backend status, please try again later")
		return
	}
	if result.State == "running" {
		SendT(b, m.Chat.ID, "There is still a test running, please wait for all works done.")
		return
	}
	nodes, err := speedtest.ReadSubscriptions(speedtest.GetHost(), url)
	if err != nil {
		SendT(b, m.Chat.ID, err.Error())
		return
	}
	if len(include) != 0 {
		nodes = speedtest.IncludeRemarks(nodes, include)
	}
	if len(exclude) != 0 {
		nodes = speedtest.ExcludeRemarks(nodes, exclude)
	}
	cfg := speedtest.NewStartConfigs(method, mode, nodes)
	go speedtest.StartTest(speedtest.GetHost(), cfg, make(chan string, 1))
	SendT(b, m.Chat.ID, "Test started, you can use /result to check latest result.")
}

func parseMsgText(s string) map[string]string {
	return ArgsParser.Parser(CfgFlags, s)
}

// cmd /run_url
func cmdStartTestWithURL(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 3 {
		SendP(b, m.Chat.ID, "Require subscriptions url.\n"+
			"Use case: <code>/run_url -u https://example.com -M TCP_PING -m ST_ASYNC</code>\n(all in upper case)", "HTML")
		return
	}
	args := parseMsgText(m.Text)
	startTestWithURL(b, m, args["-u"], args["-m"], args["-M"], []string{}, []string{})
}

// cmd /list_subs
func cmdListSubs(b *B, m *M) {
	subsFile := config.GetSubsFile()
	keys := subsFile.Section("").KeyStrings()
	if len(keys) == 0 {
		SendT(b, m.Chat.ID, "There is no subscriptions url in storage")
		return
	}
	var text string = "<b>Your subscriptions</b>:\n"
	for _, k := range keys {
		text += fmt.Sprintf("* <a href=\"%s\">%s</a>\n", subsFile.Section("").Key(k).String(), k)
	}
	SendP(b, m.Chat.ID, text, "html")
}

var Def *DefaultConfig = &DefaultConfig{
	Mode:     "TCP_PING",
	Method:   "ST_ASYNC",
	Interval: 300,
}

// cmd /set_default
func cmdSelectDefaultSub(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 2 {
		SendT(b, m.Chat.ID, "Require one arguments. \nUse case: /set_default xxx")
		return
	}

	def := strings.Fields(m.Text)[1]
	subsFile := config.GetSubsFile()
	if !subsFile.Section("").HasKey(def) {
		SendT(b, m.Chat.ID, "Remarks not found.")
		return
	}
	Def.Remarks = def
	sub := subsFile.Section("").Key(Def.Remarks).String()
	Def.Url = sub
	SendP(b, m.Chat.ID, fmt.Sprintf("Default has set to <a href=\"%s\">%s</a>", Def.Url, Def.Remarks), "HTML")
}

// cmd /set_chat
func cmdSetDefaultChat(b *B, m *M) {
	if len(strings.Fields(m.Text)) < 2 {
		SendT(b, m.Chat.ID, "Send me a char room id")
		return
	}
	val, err := strconv.ParseInt(strings.Fields(m.Text)[1], 10, 64)
	if err != nil {
		SendT(b, m.Chat.ID, err.Error())
		return
	}
	Def.Chat = val
	SendT(b, m.Chat.ID, fmt.Sprint("Default chat has set to ", Def.Chat))
}

// cmd /set_def_mode
func cmdSetDefaultModeAndMethod(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 3 {
		SendT(b, m.Chat.ID, "Require mode or method.\n"+
			"Use case:/set_def_mode -M TCP_PING -m ST_ASYNC (all in upper case)\n")
		return
	}
	args := parseMsgText(m.Text)
	Def.Mode = args["-M"]
	Def.Method = args["-m"]
	SendT(b, m.Chat.ID, "Default test mode now is "+Def.Mode+"\nDefault test method now is "+Def.Method)
}

// cmd /run_def
func cmdRunDefault(b *B, m *M) {
	startTestWithURL(b, m, Def.Url, Def.Method, Def.Mode, Def.Include, Def.Exclude)
}

// cmd /schedule
func cmdSchedule(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 2 {
		SendP(b, m.Chat.ID, "Require parameters like: <code>start/stop/status</code>\n"+
			"Use case: /schedule start", "HTML")
		return
	}
	arg := strings.Fields(m.Text)[1]
	task := NewJob()
	switch arg {
	case "start":
		if Def.Chat == 0 || Def.Url == "" || Def.Remarks == "" {
			SendT(b, m.Chat.ID, "You don't set up default config yet. Please use /set_default to set your config. Also you can check out /show_def for your current default setting.")
			return
		}

		if task.started {
			SendT(b, m.Chat.ID, "Schedule jobs has started")
			return
		}
		task.start(b)
		SendT(b, m.Chat.ID, "Jobs started")
	case "stop":
		task.Stop(0)
		SendT(b, m.Chat.ID, "Schedule jobs has stopped, but you should checkout backend for it's status.")
	case "status":
		if task.started {
			SendT(b, m.Chat.ID, "jobs running.")
			return
		}
		SendT(b, m.Chat.ID, "There is no jobs running in the background.")
	default:
		SendT(b, m.Chat.ID, "Unknown parameter.")
	}
}

// cmd /set_interval
func cmdSetInterval(b *B, m *M) {
	if len(strings.Fields(m.Text)) < 2 {
		SendT(b, m.Chat.ID, "Seconds are require as parameters\n"+
			"Use case: /set_interval 1\n"+
			"This will let the schedule task to start every 1 seconds.\n"+
			"But because of the python backend, too frequent request will cause performance problem. We recommended you use 300 second as parameter or more.")
		return
	}
	arg := strings.Fields(m.Text)[1]
	intArg, err := strconv.Atoi(arg)
	if err != nil {
		SendT(b, m.Chat.ID, "Unexpected value: "+arg)
		return
	}
	SetInterval(intArg)
	SendT(b, m.Chat.ID, "Interval has set to "+arg+"s")
}

// cmd /set_exin
func cmdSetDefaultExcludeOrInclude(b *B, m *M) {
	if args := strings.Fields(m.Text); len(args) > 2 {
		method := args[1]
		keys := args[2:]
		if method == "exclude" {
			Def.Exclude = make([]string, len(keys))
			copy(Def.Exclude, keys)
			SendT(b, m.Chat.ID, fmt.Sprintf("Default exclude has set to: %v", Def.Exclude))
			return
		}
		if method == "include" {
			Def.Include = make([]string, len(keys))
			copy(Def.Include, keys)
			SendT(b, m.Chat.ID, fmt.Sprintf("Default include has set to: %v", Def.Include))
			return
		}
		SendT(b, m.Chat.ID, "Unknown method.")
		return
	}
	SendP(b, m.Chat.ID, "Usage: /set_exin [exclude/include] keyword1 keyword2\n\n"+
		"Use case: <code>/set_exin exclude 官网 剩余流量 台 香港</code>\n"+
		"(Fuzz match is supported)", "html")
}

// cmd /show_def
func cmdShowDefault(b *B, m *M) {
	SendT(b, m.Chat.ID, fmt.Sprintf("%+v", Def))
}

// cmd /add_admin
func cmdAddAdmin(b *B, m *M) {
	if m.ReplyToMessage == nil {
		SendT(b, m.Chat.ID, "Please reply to a user.")
		return
	}
	NewName := m.ReplyToMessage.From.UserName
	if NewName == "" {
		SendT(b, m.Chat.ID, "Please set up username")
		return
	}
	NewID := m.ReplyToMessage.From.ID
	NewUser := database.Admin{
		UID:  int64(NewID),
		Name: NewName,
	}
	admins = append(admins, NewUser)
	err := database.AddAdmin(database.NewDB(), NewUser)
	if err != nil {
		SendT(b, m.Chat.ID, "Error occur when create new admin. "+err.Error())
		return
	}
	SendP(b, m.Chat.ID, fmt.Sprintf("New admin <a href=\"tg://user?id=%d\">%s</a> has set up", NewUser.UID, NewUser.Name), "HTML")
}
