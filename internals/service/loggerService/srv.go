package loggerService

import (
	"test-va/internals/entity/loggerEntity"
)

type LogSrv interface {
	Info(arg0 any, args ...any)
	Debug(arg0 any, args ...any)
	Warning(arg0 any, args ...any)
	Error(arg0 any, args ...any)
	Fatal(arg0 any, args ...any)
	Audit(record *loggerEntity.AuditLog)
}
