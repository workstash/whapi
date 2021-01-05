package logger

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/workstash/whapi/config"
)

var logger = new()
var file *os.File

//New create a new logger
func new() *log.Logger {
	var err error
	file, err = os.OpenFile(fmt.Sprintf("%s/%s.log", config.Main.LoggerFile, time.Now().Format("2006-01-02")),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	return log.New(file, "", log.LstdFlags)
}

//Close the logger file
func Close() {
	file.Close()
}

//Println similar to log.Println()
func Println(v ...interface{}) {
	logger.Println(v...)
}

//Printf similar to log.Printf()
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

//Fatal similar to log.Fatal()
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

//Fatalf similar to log.Fatalf()
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}
