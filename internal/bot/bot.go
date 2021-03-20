package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/module/config"
	"go-speedtest-bot/module/controller"
	"go-speedtest-bot/module/pastebin"
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
		log.Printf("send %q: %v", text[:10]+"...", err)
	}
}

// SendP send parsed text message
func SendP(cid int64, text string, format string) {
	msg := tgbotapi.NewMessage(cid, text)
	msg.ParseMode = format
	_, err := defBot.Send(msg)
	if err != nil {
		log.Printf("send %q: %v", text[:10]+"...", err)
	}
}

// SendTF send text with formatted content
func SendTF(cid int64, content string, args ...interface{}) {
	SendT(cid, fmt.Sprintf(content, args))
}

// SendErr send error message to given chat id
// In addition, it will print the error to log.
func SendErr(cid int64, err error) {
	log.Println(err)
	SendT(cid, err.Error())
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

func getArgs(m *M) []string {
	var strs = strings.Fields(m.Text)
	if len(strs) < 2 {
		return nil
	}
	return strs[1:]
}

// cmd /start
func cmdStart(b *B, m *M) {
	text := "Here is a bot who can help you manage all your proxy."
	SendT(m.Chat.ID, text)
}

// cmd /ping
func cmdPing(m *M) {
	args := getArgs(m)

	if args == nil {
		SendT(m.Chat.ID, "usage: /ping <runner-name>")
		return
	}

	err := controller.PingFromName(m.From.ID, args[0])
	if err != nil {
		SendErr(m.Chat.ID, err)
		return
	}
	SendTF(m.Chat.ID, "%s connected successfully", args[0])
}

// cmd /status
func cmdStatus(m *M) {
	args := getArgs(m)
	if args == nil {
		SendT(m.Chat.ID, "Usage: /status <runner-name>")
		return
	}
	status, err := controller.GetStatus(m.From.ID, args[0])
	if err != nil {
		SendErr(m.Chat.ID, err)
		return
	}
	SendT(m.Chat.ID, status)
}

// cmd /listrunner
func cmdListRunner(m *M) {
	var text = controller.ListRunner(m.From.ID)
	SendT(m.Chat.ID, text)
}

// cmd /read_sub
func cmdReadSub(m *M) {
	args := getArgs(m)
	if args == nil || len(args) < 2 {
		SendT(m.Chat.ID, "Usage: /read_sub <sub> <runner>")
		return
	}
	url := args[0]
	name := args[1]
	resp, err := controller.ReadSubscriptions(
		m.From.ID,
		url,
		name,
	)
	if err != nil {
		SendErr(m.Chat.ID, err)
		return
	}
	SendT(m.Chat.ID, resp)
}

// cmd /result
func cmdResult(m *M) {
	args := getArgs(m)
	if args == nil {
		SendT(m.Chat.ID, "Usage: /result <runner-name>")
		return
	}

	text, err := controller.Result(m.From.ID, args[0])
	if err != nil {
		SendErr(m.Chat.ID, err)
		return
	}

	if config.PasteBinEnabled() {
		resp, err := pastebin.PasteWithExpiry(
			config.GetPasteBinKey(), "nodes test result", &text, "10m",
		)
		if err != nil {
			SendErr(m.Chat.ID, err)
			return
		}
		SendTF(m.Chat.ID, "result too long, view result at %s", resp)
	}

	SendT(m.Chat.ID, text)
}

// cmd /run
func cmdStartTestWithURL(m *M) {
	args := getArgs(m)
	if args == nil {
		SendT(m.Chat.ID, "Usage: /run <runner> url=<url> method=<method> mode=<mode>")
		return
	}
	SendT(m.Chat.ID, "function under maintain")
}

// cmd /list_subs
func cmdListSubs(b *B, m *M) {
	var text string = "<b>Your subscriptions</b>:\n"
	dcs := config.GetAllDefaultConfig()
	for _, dc := range dcs {
		for _, admin := range dc.Admins {
			if m.From.ID == admin {
				text += fmt.Sprintf(`* <a href="%s">%s</a>\n`, dc.Link, dc.Name)
				break
			}
		}
	}
	SendP(m.Chat.ID, text, "html")
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
