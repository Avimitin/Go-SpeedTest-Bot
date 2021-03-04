package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/internal/ArgsParser"
	"go-speedtest-bot/internal/config"
	"go-speedtest-bot/internal/speedtest"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

// NewBot return a bot instance
func NewBot() *B {
	bot, err := tgbotapi.NewBotAPI(config.GetToken())
	if err != nil {
		log.Printf("initialize bot: %v", err)
		os.Exit(1)
	}
	return bot
}

// SendT send text message
func SendT(bot *B, cid int64, text string) {
	msg := tgbotapi.NewMessage(cid, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("send %q: %v", text[:10]+"...", err)
	}
}

// SendP send parsed text message
func SendP(bot *B, cid int64, text string, format string) {
	msg := tgbotapi.NewMessage(cid, text)
	msg.ParseMode = format
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("send %q: %v", text[:10]+"...", err)
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

	admins := NewAdmin()
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
		if admins.Auth(update.Message.From.ID) {
			go CMDHandler(bot, update.Message)
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
	runners := config.GetAllRunner()
	for _, r := range runners {
		connected := speedtest.Ping(*r)
		var text string
		if connected {
			text = "Connect to " + r.Name + " successfully"
		} else {
			text = "Unable to connect to the " + r.Name
		}
		SendT(b, m.Chat.ID, text)
	}
}

// cmd /status
func cmdStatus(b *B, m *M) {
	args := strings.Fields(m.Text)
	if len(args) < 2 {
		SendT(b, m.Chat.ID, "Usage: /status <runner-name>")
		return
	}
	runner := config.GetRunner(args[1])
	result, err := speedtest.GetStatus(*runner)
	if err != nil {
		SendT(b, m.Chat.ID, fmt.Sprint(err))
		return
	}
	SendT(b, m.Chat.ID, "Status: "+result.State)
}

// cmd /read_sub
func cmdReadSub(b *B, m *M) {
	args := strings.Fields(m.Text)
	if len(args) < 3 {
		SendT(b, m.Chat.ID, "Usage:\n/read_sub <sub> <runner>")
		return
	}
	url := args[1]
	runnername := args[2]
	runner := config.GetRunner(runnername)
	subResps, err := speedtest.ReadSubscriptions(*runner, url)
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
	args := strings.Fields(m.Text)
	if len(args) < 2 {
		SendT(b, m.Chat.ID, "Usage: /result <runner-name>")
		return
	}
	runner := config.GetRunner(args[1])
	result, err := speedtest.GetResult(*runner)
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
	args := strings.Fields(m.Text)
	if len(args) < 2 {
		SendT(b, m.Chat.ID, "Usage: /run_url <runner-name>")
		return
	}
	runner := config.GetRunner(args[1])
	result, err := speedtest.GetStatus(*runner)
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
	nodes, err := speedtest.ReadSubscriptions(*runner, url)
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
	go speedtest.StartTest(*runner, cfg, make(chan string, 1))
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
	var text string = "<b>Your subscriptions</b>:\n"
	defaultconfigs := config.GetAllDefaultConfig()
	for _, dc := range defaultconfigs {
		for _, admin := range dc.Admins {
			if m.From.ID == admin {
				text += fmt.Sprintf(`* <a href="%s">%s</a>\n`, dc.Link, dc.Name)
				break
			}
		}
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

	defname := strings.Fields(m.Text)[1]
	subsFile := config.GetDefaultConfig(defname)
	if subsFile == nil {
		SendT(b, m.Chat.ID, "Remarks not found.")
		return
	}
	for _, admin := range subsFile.Admins {
		if admin == m.From.ID {
			Def.Remarks = subsFile.Name
			Def.Url = subsFile.Link
			Def.Chat = subsFile.Chat
			SendP(b, m.Chat.ID, fmt.Sprintf("Default has set to <a href=\"%s\">%s</a>", Def.Url, Def.Remarks), "HTML")
			return
		}
	}
	SendT(b, m.Chat.ID, "you can't access the config.")
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

var (
	task = NewJob()
)

// cmd /schedule
func cmdSchedule(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 2 {
		SendP(b, m.Chat.ID, "Require parameters like: <code>start/stop/status</code>\n"+
			"Use case: /schedule start", "HTML")
		return
	}
	arg := strings.Fields(m.Text)[1]
	switch arg {
	case "start":
		if Def.Chat == 0 || Def.Url == "" || Def.Remarks == "" {
			SendT(b, m.Chat.ID, "You don't set up default config yet. Please use /set_default to set your config. Also you can check out /show_def for your current default setting.")
			return
		}

		if atomic.LoadInt32(&task.status) == RUNNING {
			SendT(b, m.Chat.ID, "There is a job running in background now.")
			return
		}
		task.start(b)
		SendT(b, m.Chat.ID, "Jobs started")
	case "stop":
		task.Stop(0)
		SendT(b, m.Chat.ID, "Schedule jobs has stopped, but you should checkout backend for it's status.")
	case "status":
		if atomic.LoadInt32(&task.status) == RUNNING {
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
