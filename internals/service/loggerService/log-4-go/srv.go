package log_4_go

import (
	"encoding/json"
	"fmt"
	logs "github.com/jeanphorn/log4go"
	"os"
	"runtime"
	"test-va/internals/entity/loggerEntity"
	"test-va/internals/service/loggerService"
)

type logSrv struct {
	logger *logs.Filter
}

func (l logSrv) Info(arg0 any, args ...any) {
	l.logger.Log(logs.INFO, getSource(), fmt.Sprintf(arg0.(string), args...))
}

func (l logSrv) Debug(arg0 any, args ...any) {
	l.logger.Log(logs.DEBUG, getSource(), fmt.Sprintf(arg0.(string), args...))
}

func (l logSrv) Warning(arg0 any, args ...any) {
	l.logger.Log(logs.WARNING, getSource(), fmt.Sprintf(arg0.(string), args...))
}

func (l logSrv) Error(arg0 any, args ...any) {
	l.logger.Log(logs.ERROR, getSource(), fmt.Sprintf(arg0.(string), args...))
}

func (l logSrv) Fatal(arg0 any, args ...any) {
	l.logger.Log(logs.CRITICAL, getSource(), fmt.Sprintf(arg0.(string), args...))
	l.logger.Close()
	os.Exit(1)
}

func (l logSrv) Audit(record *loggerEntity.AuditLog) {
	js, _ := json.Marshal(record)
	l.logger.Log(logs.INFO, getSource(), string(js))
}

func NewLogger() loggerService.LogSrv {
	//folder := "./logs"
	//logSettingsPath := "log.json"
	//appDir, err := os.Getwd()
	//if err != nil {
	//	fmt.Printf("Could not load log location >> ", err)
	//}
	//fmt.Println(appDir)
	//
	//_ = os.Mkdir(folder, os.ModePerm)
	//
	//logs.LoadConfiguration(appDir + string(os.PathSeparator) + logSettingsPath)

	return &logSrv{
		logger: logs.LOGGER("fileLogs"),
	}
}

func getSource() (source string) {
	if pc, _, line, ok := runtime.Caller(2); ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
	}
	return
}
