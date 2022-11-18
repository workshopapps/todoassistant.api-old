package reminderService

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

type ReminderSrv interface {
	SetReminder(dueDate string) error
}

type reminderSrv struct {
	*gocron.Scheduler
}

func (r *reminderSrv) SetReminder(dueDate string) error {

	// get string of date and convert it to Time.Time
	dDate, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		return err
	}

	// find time till time is expired

	duration := time.Until(dDate)

	// convert to minutes
	minutes := duration.Minutes()
	ss := fmt.Sprintf("%vs", minutes)

	r.Every(ss).Do(func() {
		log.Println("Doing... set task status to expired")
		// send mail / notification implementation
	})

	r.LimitRunsTo(1)
	r.StartAsync()
	return nil
}

func NewReminderSrv(scheduler *gocron.Scheduler) ReminderSrv {
	return &reminderSrv{Scheduler: scheduler}
}
