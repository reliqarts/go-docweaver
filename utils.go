package docweaver

import (
	"log"
	"os"
	"strings"
)

type loggerSet struct {
	Err  *log.Logger
	Info *log.Logger
	Warn *log.Logger
}
type tag string

const envKeyPrefix string = "DW_"

var loggers = GetLoggerSet()

func GetLoggerSet() *loggerSet {
	return &loggerSet{
		Err:  log.New(os.Stdout, "[Dw][err] ", log.Ldate|log.Ltime|log.Lshortfile),
		Info: log.New(os.Stdout, "[Dw][info] ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn: log.New(os.Stdout, "[Dw][warn] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func envKeyName(key string) string {
	return strings.ToUpper(envKeyPrefix + key)
}
