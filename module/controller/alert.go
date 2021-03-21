package controller

import (
	"fmt"
	"go-speedtest-bot/module/speedtest"
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

func AlertHandler(results []speedtest.ResultInfo) string {
	var text string
	for _, r := range results {
		if r.Ping < 0.0001 || r.GPing < 0.0001 {
			if !HasRecode(r.Remarks) {
				text += alertNotifyLn(r.Remarks)
			}
			AppendDiag(r.Remarks)
			continue
		}
		if HasRecode(r.Remarks) {
			DelRecord(r.Remarks)
			record := CheckRecord(r.Remarks)
			text += recoverNotifyLn(r.Remarks, record.OfflineDuration, record.Count)
		}
	}
	return text
}

func alertNotifyLn(remark string) string {
	return fmt.Sprintf("%s offline\n", remark)
}

func recoverNotifyLn(remark string, duration time.Duration, count int) string {
	return fmt.Sprintf("%s has recovered, offline nearly %v, %d times\n", remark, duration, count)
}
