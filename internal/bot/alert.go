package bot

import (
	"go-speedtest-bot/internal/speedtest"
	"time"
)

type Recode struct {
	Count           int
	Exist           bool
	Date            time.Time
	OfflineDuration time.Duration
}
type History = map[string]*Recode

var DiagHistory = make(History)

func SetAlert(enable bool) {
	alert = enable
}

func CheckDiag() History {
	return DiagHistory
}

func CheckRecord(recode string) *Recode {
	if v, ok := CheckDiag()[recode]; ok {
		return v
	}
	return nil
}

func HasRecode(recode string) bool {
	if result := CheckRecord(recode); result != nil {
		return result.Exist
	}
	return false
}

func AppendDiag(record string) {
	if val, ok := DiagHistory[record]; ok {
		// If the node has offline, add it's count
		if !val.Exist {
			val.Exist = true
			val.Count++
			val.Date = time.Now()
		}
		return
	}
	DiagHistory[record] = &Recode{1, true, time.Now(), 0}
}

func DelRecord(record string) {
	if val, ok := DiagHistory[record]; ok {
		if val.Exist == true {
			val.Exist = false
			val.OfflineDuration = time.Now().Sub(val.Date)
		}
	}
}

func AlertHandler(results []speedtest.ResultInfo, b *B) {
	for _, r := range results {
		if r.Ping < 0.0001 || r.GPing < 0.0001 {
			if !HasRecode(r.Remarks) {
				Alert(b, r.Remarks)
			}
			AppendDiag(r.Remarks)
			continue
		}
		if HasRecode(r.Remarks) {
			DelRecord(r.Remarks)
			RecoverNotice(b, r.Remarks, CheckRecord(r.Remarks).OfflineDuration.String())
		}
	}
}

func Alert(b *B, remark string) {
	SendT(b, Def.Chat, remark+" offline.")
}

func RecoverNotice(b *B, remark string, duration string) {
	SendT(b, Def.Chat, remark+" now online. Offline nearly "+duration)
}
