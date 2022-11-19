package timeSrv

import "time"

type TimeService interface {
	CurrentTime() time.Time
	TimeSince(time2 time.Time) time.Duration
	CheckFor339Format(time string) error
}

type timeStruct struct{}

func (t timeStruct) CheckFor339Format(timeStr string) error {
	_, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}
	// add check if it is greater than current time

	return nil
}

func NewTimeStruct() TimeService {
	return &timeStruct{}
}

func (t timeStruct) CurrentTime() time.Time {
	return time.Now()
}

func (t timeStruct) TimeSince(time2 time.Time) time.Duration {
	return time.Since(time2)
}
