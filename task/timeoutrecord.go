package task

import (
	"library/models/records"
	"log"
	"time"
)

func VerifyTimeoutRecord() {
	timer := time.NewTimer(time.Hour * 24)

	go handler(timer)
}

func handler(timer *time.Timer) {
	<-timer.C

	all, err := records.All()
	if err != nil {
		return
	}

	for _, v := range all {
		if v.EndDate.Unix() < time.Now().Unix() {
			newDate := v.EndDate.Add(time.Hour * 24)
			err = records.UpdateEndDate(v.ID, &newDate)
			if err != nil {
				log.Println("delay end date err", err)
			}
		}
	}

	VerifyTimeoutRecord()
}
