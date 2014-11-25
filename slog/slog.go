package slog

// 暂时没有加锁，注意只能在GOPROC == 1的情况下使用

import (
	"log"
//	"io"
	"os"
	"time"
	"fmt"

)

type logger struct {
	logpref string

	loghour string
	logfp *os.File
	per *log.Logger

}

func (self *logger) setOutput() {
	hour := time.Now().Format("2006-01-02-15")
	//log.Println("setoutput", hour)
	if self.logpref == "" && self.loghour == "" {
		self.per = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		self.loghour = hour
		//log.Println("setoutput", "std", hour)
	}

	if self.logpref != "" && self.loghour != hour {
		logFile := fmt.Sprintf("%s.%s.log", self.logpref, hour)
		logf, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
			return
		}

		//log.Println("setoutput", "pref", self.logpref, hour)

		self.per = log.New(logf, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		if self.logfp != nil {
			self.logfp.Close()
		}
		self.logfp = logf
		self.loghour = hour
	}


}

func (self *logger) Printf(format string, v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}
	self.per.Printf(format, v...)
}

func (self *logger) Panicf(format string, v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}

	self.per.Panicf(format, v...)
}


func (self *logger) Println(v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}

	self.per.Println(v...)
}

func (self *logger) Panicln(v ...interface{}) {
	self.setOutput()
	if self.per == nil {
		log.Println("slog nil")
		return
	}

	self.per.Panicln(v...)
}

// log 级别
const (
	LV_TRACE int = 0
	LV_DEBUG int = 1
	LV_INFO int = 2
	LV_WARN int = 3
	LV_ERROR int = 4
	LV_FATAL int = 5
	LV_PANIC int = 6


)



var (
	headTrace string
	headDebug string
	headInfo string
	headWarn string
	headError string
	headFatal string
	headPanic string

	headFmtTrace string
	headFmtDebug string
	headFmtInfo string
	headFmtWarn string
	headFmtError string
	headFmtFatal string
	headFmtPanic string

	log_level int
	lg *logger
)

func Init(logdir string, logpref string, level string) {
	if level == "TRACE" {
		log_level = LV_TRACE
	} else if level == "DEBUG" {
		log_level = LV_DEBUG
	} else if level == "INFO" {
		log_level = LV_INFO
	} else if level == "WARN" {
		log_level = LV_WARN
	} else if level == "ERROR" {
		log_level = LV_ERROR
	} else if level == "FATAL" {
		log_level = LV_FATAL
	} else if level == "PANIC" {
		log_level = LV_PANIC
	} else {
		log_level = LV_INFO
	}

	if logdir != "" {
		err := os.MkdirAll(logdir, 0777)
		if err != nil {
			log.Fatalln("slog mkdir ", logdir, " err:", err)
		}
	}

	logfile := ""
	if logdir != "" && logpref != "" {
		logfile = logdir+"/"+logpref
	}

    lg = &logger{logpref: logfile, logfp: nil, per: nil}


}

func init() {
	pid := os.Getpid()
	headTrace = fmt.Sprintf("[TRACE] [%d]", pid)
	headDebug = fmt.Sprintf("[DEBUG] [%d]", pid)
	headInfo = fmt.Sprintf("[INFO] [%d]", pid)
	headWarn = fmt.Sprintf("[WARN] [%d]", pid)
	headError = fmt.Sprintf("[ERROR] [%d]", pid)
	headFatal = fmt.Sprintf("[FATAL] [%d]", pid)
	headPanic = fmt.Sprintf("[PANIC] [%d]", pid)

	headFmtTrace = fmt.Sprintf("%s ", headTrace)
	headFmtDebug = fmt.Sprintf("%s ", headDebug)
	headFmtInfo = fmt.Sprintf("%s ", headInfo)
	headFmtWarn = fmt.Sprintf("%s ", headWarn)
	headFmtError = fmt.Sprintf("%s ", headError)
	headFmtFatal = fmt.Sprintf("%s ", headFatal)
	headFmtPanic = fmt.Sprintf("%s ", headPanic)


	Init("", "", "TRACE")
}


func Tracef(format string, v ...interface{}) {
	if LV_TRACE >= log_level {
		lg.Printf(headFmtTrace+format, v...)
	}
}

func Traceln(v ...interface{}) {
	if LV_TRACE >= log_level {
		lg.Println(append([]interface{}{headTrace}, v...)...)
	}
}


func Debugf(format string, v ...interface{}) {
	if LV_DEBUG >= log_level {
		lg.Printf(headFmtDebug+format, v...)
	}
}

func Debugln(v ...interface{}) {
	if LV_DEBUG >= log_level {
		lg.Println(append([]interface{}{headDebug}, v...)...)
	}
}


func Infof(format string, v ...interface{}) {
	if LV_INFO >= log_level {
		lg.Printf(headFmtInfo+format, v...)
	}
}

func Infoln(v ...interface{}) {
	if LV_INFO >= log_level {
		lg.Println(append([]interface{}{headInfo}, v...)...)
	}
}


func Warnf(format string, v ...interface{}) {
	if LV_WARN >= log_level {
		lg.Printf(headFmtWarn+format, v...)
	}
}

func Warnln(v ...interface{}) {
	if LV_WARN >= log_level {
		lg.Println(append([]interface{}{headWarn}, v...)...)
	}
}


func Errorf(format string, v ...interface{}) {
	if LV_ERROR >= log_level {
		lg.Printf(headFmtError+format, v...)
	}
}

func Errorln(v ...interface{}) {
	if LV_ERROR >= log_level {
		lg.Println(append([]interface{}{headError}, v...)...)
	}
}



func Fatalf(format string, v ...interface{}) {
	if LV_FATAL >= log_level {
		lg.Printf(headFmtFatal+format, v...)
	}
}


func Fatalln(v ...interface{}) {
	if LV_FATAL >= log_level {
		lg.Println(append([]interface{}{headFatal}, v...)...)
	}
}


func Panicf(format string, v ...interface{}) {
	if LV_PANIC >= log_level {
		lg.Panicf(headFmtPanic+format, v...)
	}
}


func Panicln(v ...interface{}) {
	if LV_PANIC >= log_level {
		lg.Panicln(append([]interface{}{headPanic}, v...)...)
	}
}