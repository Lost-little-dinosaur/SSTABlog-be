package logger

import (
	"SSTABlog-be/config"
	"io"
	"log"
	"os"
)

type debugDefault struct {
	Debug *log.Logger
}

func (d *debugDefault) Println(v ...interface{}) {
	if config.GetConfig().MODE == "debug" {
		d.Debug.Println(v)
	}
}

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *debugDefault
)

func init() {
	errFile, err := os.OpenFile(ERROR_LOG_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开ERROR日志文件失败！")
	}
	warningFile, err := os.OpenFile(WARING_LOG_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开WARNING日志文件失败！")
	}
	infoFile, err := os.OpenFile(INFO_LOG_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开INFO日志文件失败！")
	}
	debugFile, err := os.OpenFile(DEBUG_LOG_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开DEBUG日志文件失败！")
	}

	Info = log.New(io.MultiWriter(os.Stderr, infoFile, os.Stdout), "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stderr, warningFile, os.Stdout), "[Warning] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = &debugDefault{
		Debug: log.New(io.MultiWriter(os.Stderr, debugFile, os.Stdout), "[Debug] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
