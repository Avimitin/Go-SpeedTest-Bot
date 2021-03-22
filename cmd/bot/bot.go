package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/module/config"
	"go-speedtest-bot/module/controller"
	"go-speedtest-bot/module/pastebin"
	"log"
	"os"
	"strings"
)

var (
	defBot *B
	rc     = new(RunnerComm)
)

// NewBot return a bot instance
func NewBot() *B {
	config.LoadConfig()

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
		if len(text) > 10 {
			text = text[:10] + "..."
		}
		log.Printf("send %q : %v", text, err)
	}
}

// SendP send parsed text message
func SendP(cid int64, text string, format string) {
	msg := tgbotapi.NewMessage(cid, text)
	msg.ParseMode = format
	_, err := defBot.Send(msg)
	if err != nil {
		if len(text) > 10 {
			text = text[:10] + "..."
		}
		log.Printf("send %q: %v", text[:10]+"...", err)
	}
}

// SendTF send text with formatted content
func SendTF(cid int64, content string, args ...interface{}) {
	SendT(cid, fmt.Sprintf(content, args...))
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
func Listen(clean bool) {
	defBot = NewBot()
	if os.Getenv("debug_bot") == "true" {
		defBot.Debug = true
	}
	log.Println("Authorized on account", defBot.Self.UserName)

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

		go CMDHandler(update.Message)
	}
}

// CMDHandler handle all the command
func CMDHandler(msg *M) {
	if msg.IsCommand() {
		cmd := msg.Command()
		log.Printf("[%s]%s", msg.From.FirstName, cmd)
		if cmd, ok := Commands[cmd]; ok {
			cmd(msg)
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
func cmdStart(m *M) {
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
	SendT(m.Chat.ID, "handling sub")
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
			config.GetPasteBinKey(), "nodes test result", &text, "10M",
		)
		if err != nil {
			SendErr(m.Chat.ID, err)
			return
		}
		SendTF(m.Chat.ID, "result too long, view result at %s", resp)
		return
	}

	SendT(m.Chat.ID, text)
}

// cmd /run
func cmdStartTestWithURL(m *M) {
	args := getArgs(m)
	if args == nil || len(args) < 2 {
		SendT(m.Chat.ID, "Usage: /run <runner> <json-args>")
		return
	}
	var rname = args[0]
	var jargs = args[1]

	var parsedArgs *SpeedTestArguments
	err := json.Unmarshal([]byte(jargs), &parsedArgs)
	if err != nil {
		SendErr(m.Chat.ID, fmt.Errorf("argument not valid: %v", err))
		return
	}

	var errCh = make(chan error)
	err = controller.Run(
		m.From.ID,
		rname,
		parsedArgs.Subs,
		parsedArgs.Method,
		parsedArgs.Mode,
		parsedArgs.Include,
		parsedArgs.Exclude,
		errCh,
	)
	if err != nil {
		SendErr(m.Chat.ID, err)
		return
	}
	go func() {
		err := <-errCh
		if err != nil {
			SendErr(m.Chat.ID, err)
			return
		}
		SendT(m.Chat.ID, rname+" finished a job")
	}()

	SendT(m.Chat.ID, "request success, checkout result by /result")
}

// cmd /list_subs
func cmdListSubs(m *M) {
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
func cmdSchedule(m *M) {
	args := getArgs(m)
	if args == nil || len(args) < 2 {
		SendP(m.Chat.ID, "Require parameters like: <code>start/stop/status</code>\n"+
			"Use case: /schedule start <code>{{CONFIG_NAME}}</code>", "HTML")
		return
	}

	def := config.GetDefaultConfig(args[1])
	if def == nil {
		SendT(m.Chat.ID, "config specific not found.")
		return
	}
	if !def.HasAccess(m.From.ID) {
		SendT(m.Chat.ID, "you don't have access to this config")
	}

	switch args[0] {
	case "start":
		if rc.Exist(def.Name) {
			SendT(m.Chat.ID, "there is another schedule jobs running.")
		}

		c := controller.NewComm()
		err := controller.Schedule(m.From.ID, def, c)
		if err != nil {
			SendErr(m.Chat.ID, err)
			return
		}
		// if no error register the channel
		rc.Register(def.Name, c)

		go ScheduleJobsNotify(m.Chat.ID, c)
		SendT(m.Chat.ID, "Jobs started")
	case "stop":
		c := rc.C(def.Name)
		if c == nil {
			SendT(m.Chat.ID, "runner specific is not running schedule jobs now")
			return
		}

		c.Sig <- 0
		rc.UnRegister(def.Name)

		SendT(m.Chat.ID, "jobs has exit at local, but the remote host maybe still running"+
			" jobs, please check out the backend for its status")
	case "status":
		if !rc.Exist(def.Name) {
			SendT(m.Chat.ID, "runner is pending.")
			return
		}
		SendT(m.Chat.ID, "runner are handling schedule test jobs.")
	default:
		SendT(m.Chat.ID, "Unknown parameter.")
	}
}

func cmdListFailed(m *M) {
	var header = "current offline nodes:\n\n"
	var text string

	diag := controller.CheckDiag()
	for k, v := range diag {
		if v.Exist {
			text += k
		}
	}

	if len(text) == 0 {
		SendT(m.Chat.ID, "no nodes offline")
		return
	}

	SendT(m.Chat.ID, header+text)
}
