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

func SendP(bot *B, cid int64, text string, format string) {
	msg := tgbotapi.NewMessage(cid, text)
	msg.ParseMode = format
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("SendError", err)
	}
}

func CMDHandler(bot *B, msg *M) {
	if msg.IsCommand() {
		if cmd, ok := Commands[msg.Command()]; ok {
			cmd(bot, msg)
		}
	}
}

func cmdStart(b *B, m *M) {
	text := "Here is a bot who can help you manage all your proxy."
	SendT(b, m.Chat.ID, text)
}

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

func cmdStatus(b *B, m *M) {
	result, err := speedtest.GetStatus(speedtest.GetHost())
	if err != nil {
		SendT(b, m.Chat.ID, fmt.Sprint(err))
		return
	}
	SendT(b, m.Chat.ID, "Status: "+result.State)
}

func cmdReadSub(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) != 3 {
		SendT(b, m.Chat.ID, "Use case(Only single link is supported):\n/read_sub https://xxx.com")
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
			"%s: | ls: %.2f%% | lp: %.2f ms | gp: %.2f ms", res.Remarks, res.Loss*100, res.Ping*1000, res.GPing*1000)
	}
	return text
}

func cmdResult(b *B, m *M) {
	result, err := speedtest.GetResult(speedtest.GetHost())
	if err != nil {
		SendT(b, m.Chat.ID, fmt.Sprint(err))
		return
	}
	if fresult := formatResult(result); len(fresult) != 0 {
		SendT(b, m.Chat.ID, fresult)
	}
	SendT(b, m.Chat.ID, "No result yet")
}

func startTestWithURL(b *B, m *M, url string, method string, mode string) {
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
	cfg := speedtest.NewStartConfigs(method, mode, nodes)
	speedtest.StartTest(speedtest.GetHost(), cfg)
	SendT(b, m.Chat.ID, "Test started, you can use /result to check latest result.")
}

func parseMsgText(b *B, m *M) map[string]string {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) != 3 {
		SendT(b, m.Chat.ID, "Require subscriptions url.\n"+
			"Use case:/run_url -u https://example.com -M TCP_PING -m ST_ASYNC (all in upper case)\n")
		return nil
	}
	return ArgsParser.Parser(CfgFlags, m.Text)
}

func cmdStartTestWithURL(b *B, m *M) {
	args := parseMsgText(b, m)
	startTestWithURL(b, m, args["-u"], args["-M"], args["-m"])
}

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
	Mode:   "TCP_PING",
	Method: "ST_ASYNC",
}

func cmdSelectDefaultSub(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) != 3 {
		SendT(b, m.Chat.ID, "Require one arguments. \n Use case: /default xxx")
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
	SendT(b, m.Chat.ID, "Default has set to "+Def.Remarks+"\n"+"url: "+sub)
}

func cmdSetDefaultModeAndMethod(b *B, m *M) {
	if len(m.Text)-1 == len(m.Command()) || len(strings.Fields(m.Text)) != 3 {
		SendT(b, m.Chat.ID, "Require mode or method.\n"+
			"Use case:/set_def_mode -M TCP_PING -m ST_ASYNC (all in upper case)\n")
		return
	}
	args := parseMsgText(b, m)
	Def.Mode = args["-M"]
	Def.Method = args["-m"]
	SendT(b, m.Chat.ID, "Default test mode now is "+Def.Mode+"\nDefault test method now is "+Def.Method)
}

func cmdRunDefault(b *B, m *M) {
	startTestWithURL(b, m, Def.Url, Def.Method, Def.Mode)
}
