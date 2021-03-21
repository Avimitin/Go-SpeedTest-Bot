package bot

import (
	"go-speedtest-bot/module/controller"
	"log"
)

func ScheduleJobsNotify(cid int64, c *controller.Comm) {
	for {
		select {
		case e := <-c.ErrCh:
			SendErr(cid, e)
		case l := <-c.LogCh:
			log.Println(l)
		case a := <-c.Alert:
			SendT(cid, *a)
		}
	}
}
