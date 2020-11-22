package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/internal/speedtest"
	"log"
	"testing"
	"time"
)

func TestLoadConf(t *testing.T) {
	if loadConf().token == "" {
		t.Fail()
	}
	fmt.Println(loadConf().token)
}

func TestNewBot(t *testing.T) {
	bot := NewBot()
	me, err := bot.GetMe()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println("Auth on: " + me.UserName)
}

func TestSendT(t *testing.T) {
	SendT(NewBot(), 649191333, "Test message from TestSendT function")
}

func TestSendP(t *testing.T) {
	SendP(NewBot(), 649191333, "*Test message* from `TestSendP` function", "markdownv2")
}

func TestLaunch(t *testing.T) {
	Launch(true, true, true)
}

func TestCMDHandler(t *testing.T) {
	Commands := CMD{
		"test": func(b *B, m *M) {
			fmt.Println(b.Self.UserName, m.Text)
		},
	}
	msg := &tgbotapi.Message{
		Text: "/test message",
		Entities: &[]tgbotapi.MessageEntity{
			{Offset: 0, Type: "bot_command", Length: 5},
		},
	}
	if msg.IsCommand() {
		if cmd, ok := Commands[msg.Command()]; ok {
			cmd(NewBot(), msg)
		} else {
			log.Println("command not in list.")
			t.Fail()
		}
	} else {
		log.Println("msg is not a command.")
		t.Fail()
	}
}

func NewMsg(msgText string, cmdLen int) *M {
	return &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 649191333},
		Text:     msgText,
		Entities: &[]tgbotapi.MessageEntity{{Offset: 0, Type: "bot_command", Length: cmdLen}},
	}
}

func TestCMDStart(t *testing.T) {
	cmdStart(NewBot(), NewMsg("/start", 6))
}

func TestCMDPing(t *testing.T) {
	cmdPing(NewBot(), NewMsg("/ping", 5))
}

func TestCMDStatus(t *testing.T) {
	cmdStatus(NewBot(), NewMsg("/status", 7))
}

func TestCMDReadSub(t *testing.T) {
	cmdReadSub(NewBot(), NewMsg("/read_sub ", 9))
}

func TestCMDResult(t *testing.T) {
	cmdResult(NewBot(), NewMsg("/result", 7))
}

func TestCmdList(t *testing.T) {
	cmd := "/list_subs"
	cmdListSubs(NewBot(), NewMsg(cmd, len(cmd)))
}

func TestParse(t *testing.T) {
	args := parseMsgText("/run_url -u https://oxygenproxy.com/auth/register")
	if args["-u"] != "https://oxygenproxy.com/auth/register" {
		t.Fail()
	}
}

func TestCMDRunURL(t *testing.T) {
	cmd := "/run_url"
	cmdStartTestWithURL(NewBot(), NewMsg(cmd+" -u https://oxygenproxy.com/auth/register", len(cmd)))
}

func TestCMDSelectDef(t *testing.T) {
	cmd := "/set_default"
	cmdSelectDefaultSub(NewBot(), NewMsg(cmd+" OxygenProxy", len(cmd)))
}

func TestSetDefMode(t *testing.T) {
	cmd := "/set_def_mode"
	cmdSetDefaultModeAndMethod(NewBot(), NewMsg(cmd+" -M TCPPPP", len(cmd)))
}

func TestRunDef(t *testing.T) {
	Def.Url = ""
	Def.Include = []string{"香港"}
	cmdRunDefault(NewBot(), NewMsg("/run_def", 8))
}

func TestSchedule(t *testing.T) {
	Def.Interval = 10
	Def.Chat = 649191333
	Def.Include = []string{"剩余"}
	start(NewBot())
}

func TestCMDSchedule(t *testing.T) {
	cmd := "/schedule"
	cmdSchedule(NewBot(), NewMsg(cmd+" status", len(cmd)))
}

func TestSetDefaultExIn(t *testing.T) {
	cmd := "/set_exin"
	cmdSetDefaultExcludeOrInclude(NewBot(), NewMsg(cmd, len(cmd)))
}

func TestSetAlert(t *testing.T) {
	SetAlert(true)
	if !alert {
		t.Fail()
	}
}

func TestCheckDiag(t *testing.T) {
	get := CheckDiag()
	for i := range get {
		if get[i] != DiagHistory[i] {
			t.Fail()
		}
	}
}

func TestAppendDiag(t *testing.T) {
	AppendDiag("a")
	if CheckDiag()["a"].Count != 1 {
		t.Errorf("Got %d want 1", CheckDiag()["a"].Count)
	}
	AppendDiag("b")
	if CheckDiag()["b"].Count != 1 {
		t.Errorf("Got %d want 1", CheckDiag()["b"].Count)
	}
	AppendDiag("a")
	if CheckDiag()["a"].Count == 2 {
		t.Errorf("Got %d want 1", CheckDiag()["a"].Count)
	}
	DelRecord("a")
	AppendDiag("a")
	e := CheckDiag()["a"].Exist
	d := CheckRecord("a").Date
	c := CheckDiag()["a"].Count
	if !e && c != 2 && d != time.Now() {
		t.Errorf("Want a exist and count 2 got a %v and count %d and time not right", e, c)
	}
}

func TestDelRecord(t *testing.T) {
	for i := 3; i > 0; i-- {
		AppendDiag("a")
	}
	DelRecord("a")
	if CheckDiag()["a"].Exist {
		t.Errorf("Got exist want not exist.")
	}
	AppendDiag("a")
	if count := CheckDiag()["a"].Count; count != 2 {
		t.Errorf("Got %d want 2", count)
	}
}

func TestAlertHandler(t *testing.T) {
	testResults := []speedtest.ResultInfo{
		{
			Remarks: "HK-JP01",
			Ping:    0.00,
			GPing:   0.00,
		},
		{
			Remarks: "SZ-HK01",
			Ping:    0.10,
			GPing:   0.45,
		},
	}
	b := NewBot()
	Def.Chat = 649191333
	AlertHandler(testResults, b)
	if !HasRecode("HK-JP01") {
		t.Errorf("Want nodes exist but got null")
	}
	if HasRecode("SZ-HK01") {
		t.Errorf("Unwanted node exist in history.")
	}
	testResults = append(testResults, speedtest.ResultInfo{
		Remarks: "HK-JP01",
		Ping:    0.12,
		GPing:   0.45,
	})
	CheckRecord("HK-JP01").Date = time.Date(2020, 11, 12, 23, 30, 0, 0, time.Local)
	AlertHandler(testResults, b)
	if HasRecode("HK-JP01") {
		t.Errorf("Want node not exist but still got it.")
	}
}

func TestCheckRecord(t *testing.T) {
	AppendDiag("a")
	if !CheckRecord("a").Exist {
		t.Errorf("Want a exist but not exist.")
	}
}

func TestHasRecode(t *testing.T) {
	AppendDiag("a")
	if !HasRecode("a") {
		t.Errorf("Want a exist but got none.")
	}
	if HasRecode("b") {
		t.Errorf("Unwanted record in history.")
	}
}
