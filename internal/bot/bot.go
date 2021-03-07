package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/internal/ArgsParser"
	"go-speedtest-bot/internal/config"
	"go-speedtest-bot/internal/speedtest"
	"log"
	"os"
	"strings"
)

var defBot *B

// NewBot return a bot instance
func NewBot() *B {
	bot, err := tgbotapi.NewBotAPI(config.GetToken())
	if err != nil {
		log.Fatalf("initialize bot: %v", err)
	}
	return bot
}

// SendT send text message
func SendT(cid int64, text string) {
	msg := tgbotapi.NewMessage(cid, text)
	_, err := defBot.Send(msg)
	if err != nil {
		log.Println("send %q: %v", text[:10]+"...", err)
	}
}

// SendP send parsed text message
func SendP(cid int64, text string, format string) {
	msg := tgbotapi.NewMessage(cid, text)
	msg.ParseMode = format
	_, err := defBot.Send(msg)
	if err != nil {
		log.Println("send %q: %v", text[:10]+"...", err)
	}
}

// SendTF send text with formated content
func SendTF(cid int64, content string, args ...interface{}) {
	SendT(cid, fmt.Sprintf(content, args))
}

// Listen is the robot's main methods. Calling this method will
// enter a loop and continuously listen for messages.
// If debug, bot will enable debug mode.
// If logInfo, program will print all the message.
// If clean, program will clean all the out of date message.
func Listen(debug bool, logInfo bool, clean bool) {
	defBot = NewBot()
	defBot.Debug = debug
	log.Println("Authorized on account", defBot.Self.UserName)

	admins := NewAdmin()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := defBot.GetUpdatesChan(u)
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
			go CMDHandler(defBot, update.Message)
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
	SendT(m.Chat.ID, text)
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
		SendT(m.Chat.ID, text)
	}
}

// cmd /status
func cmdStatus(b *B, m *M) {
	args := strings.Fields(m.Text)
	if len(args) < 2 {
		SendT(m.Chat.ID, "Usage: /status <runner-name>")
		return
	}
	runner := config.GetRunner(args[1])
	result, err := speedtest.GetStatus(*runner)
	if err != nil {
		SendT(m.Chat.ID, fmt.Sprint(err))
		return
	}
	SendT(m.Chat.ID, "Status: "+result.State)
}

// cmd /read_sub
func cmdReadSub(b *B, m *M) {
	args := strings.Fields(m.Text)
	if len(args) < 3 {
		SendT(m.Chat.ID, "Usage:\n/read_sub <sub> <runner>")
		return
	}
	url := args[1]
	runnername := args[2]
	runner := config.GetRunner(runnername)
	subResps, err := speedtest.ReadSubscriptions(*runner, url)
	if err != nil {
		SendT(m.Chat.ID, fmt.Sprint(err))
		return
	}
	var text string
	for _, subResp := range subResps {
		if subResp.Type != "" && subResp.Config != nil {
			text += subResp.Type + " " + subResp.Config.Remarks + "\n"
		}
	}
	SendT(m.Chat.ID, text)
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
		SendT(m.Chat.ID, "Usage: /result <runner-name>")
		return
	}
	runner := config.GetRunner(args[1])
	result, err := speedtest.GetResult(*runner)
	if err != nil {
		SendT(m.Chat.ID, fmt.Sprint(err))
		return
	}
	if fresult := formatResult(result); len(fresult) != 0 {
		SendT(m.Chat.ID, fresult)
		return
	}
	SendT(m.Chat.ID, "No result yet")
}

func startTestWithURL(b *B, m *M, url string, method string, mode string, include []string, exclude []string) {
	args := strings.Fields(m.Text)
	if len(args) < 2 {
		SendT(m.Chat.ID, "Usage: /run_url <runner-name>")
		return
	}
	runner := config.GetRunner(args[1])
	result, err := speedtest.GetStatus(*runner)
	if err != nil {
		SendT(m.Chat.ID, err.Error())
		return
	}
	if result == nil {
		SendT(m.Chat.ID, "Unable to fetch backend status, please try again later")
		return
	}
	if result.State == "running" {
		SendT(m.Chat.ID, "There is still a test running, please wait for all works done.")
		return
	}
	nodes, err := speedtest.ReadSubscriptions(*runner, url)
	if err != nil {
		SendT(m.Chat.ID, err.Error())
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
	SendT(m.Chat.ID, "Test started, you can use /result to check latest result.")
}

func parseMsgText(s string) map[string]string {
	return ArgsParser.Parser(CfgFlags, s)
}

// cmd /run_url
func cmdStartTestWithURL(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 3 {
		SendP(m.Chat.ID, "Require subscriptions url.\n"+
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
	SendP(m.Chat.ID, text, "html")
}

// cmd /set_default
func cmdSelectDefaultSub(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) < 2 {
		SendT(m.Chat.ID, "Require one arguments. \nUse case: /set_default xxx")
		return
	}

	defname := strings.Fields(m.Text)[1]
	subsFile := config.GetDefaultConfig(defname)
	if subsFile == nil {
		SendT(m.Chat.ID, "Remarks not found.")
		return
	}
}

// cmd /schedule
func cmdSchedule(b *B, m *M) {
	if len(strings.Fields(m.Text)) < 3 {
		SendP(m.Chat.ID, "Require parameters like: <code>start/stop/status</code>\n"+
			"Use case: /schedule start <CONFIG_NAME>", "HTML")
		return
	}
	args := strings.Fields(m.Text)

	subsFile := config.GetDefaultConfig(args[2])
	if subsFile == nil {
		SendT(m.Chat.ID, "config specific not found.")
		return
	}

	runner := config.GetRunner(subsFile.DefaultRunner)
	if runner == nil {
		SendT(m.Chat.ID,
			"the runner name specific in default config is not found, please check your config")
		return
	}

	switch args[1] {
	case "start":
		var haveAccess bool
		for _, admin := range subsFile.Admins {
			if admin == m.From.ID {
				haveAccess = true
				break
			}
		}
		if !haveAccess {
			SendT(m.Chat.ID, "you don't have access to this profile")
			return
		}
		if runner.IsWorking() {
			SendTF(m.Chat.ID, "runner %s is handling other work", runner.Name)
			return
		}
		go StartScheduleJobs(runner, subsFile)
	case "stop":
		runner.HangUp()
		runner.CloseChan()
		SendT(m.Chat.ID,
			"Schedule jobs has been stopped, "+
				"but backend speedtest doesn't stop"+
				"immediately, so please checkout backend status"+
				"yourself before starting new test.")
	case "status":
		if runner.IsPending() {
			SendT(m.Chat.ID, "runner is pending.")
			return
		}
		SendT(m.Chat.ID, "runner are handling speedtest work.")
	default:
		SendT(m.Chat.ID, "Unknown parameter.")
	}
}
