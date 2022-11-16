package timeSrv

import "time"

type TimeService interface {
	CurrentTime() time.Time
	TimeSince(time2 time.Time) time.Duration
}

type timeStruct struct{}

func NewTimeStruct() TimeService {
	return &timeStruct{}
}

func (t timeStruct) CurrentTime() time.Time {
	return time.Now()
}

func (t timeStruct) TimeSince(time2 time.Time) time.Duration {
	return time.Since(time2)
}
