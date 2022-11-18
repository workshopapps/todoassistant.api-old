package reminderService

import (
	"github.com/go-co-op/gocron"
	"testing"
	"time"
)

func Test_reminderSrv_SetReminder(t *testing.T) {
	gcrn := gocron.NewScheduler(time.UTC)
	srv := NewReminderSrv(gcrn)
	due := time.Now().Add(2 * time.Minute).Format(time.RFC3339)
	srv.SetReminder(due)
}
