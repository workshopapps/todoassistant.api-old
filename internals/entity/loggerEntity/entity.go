package loggerEntity

type LogRecord struct {
	Level    string // The log level
	Date     string // The time at which the log message was created (nanoseconds)
	Source   string // The message source
	Message  string // The log message
	Category string // The log group
}

// AuditLog Audit log
type AuditLog struct {
	Date           string  `json:"date"`
	Username       string  `json:"Username"`
	RequestHeader  any     `json:"request_header"`
	Request        any     `json:"request"`
	StatusCode     int     `json:"status_code"`
	ResponseHeader any     `json:"response_header"`
	Response       any     `json:"response"`
	ClientID       string  `json:"client_id"`
	Route          string  `json:"route"`
	Duration       float64 `json:"duration (seconds)"`
}
