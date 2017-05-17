package utils

import (
	"sync"
	"time"
)

var oneTime sync.Once
var logInstance *logging

const TRACE int = 0
const DEBUG int = 1
const INFO int = 2
const WARN int = 3
//const ERROR int = 4


type logging struct{

	logStream chan string
	level int
}

func GetLogging() *logging {
	oneTime.Do(func() {
		logInstance = &logging{}
		logInstance.init()
	})
	return logInstance
}

func (l *logging) init(){
		l.logStream = make(chan string)
		l.level = GetSettings().GetLevel();
		go l.logWriter()
}

func (l *logging) Trace(s string){
	if(l.level == TRACE){
		l.logStream <- "TRACE: "+s
	}
}
func (l *logging) Debug(s string){
	if(l.level <= DEBUG){
		l.logStream <- "DEBUG: "+s
	}
}
func (l *logging) Info(s string){
	if(l.level <= INFO){
		l.logStream <- "INFO: "+s
	}
}
func (l *logging) Warn(s string){
	if(l.level <= WARN){
		l.logStream <- "WARN: "+s
	}
}
func (l *logging) Error(s string){
	//Error aways gets logged
	l.logStream <- "ERROR: "+s

}

func (l *logging) logWriter(){

	for {
		statement :=<-l.logStream
		dt := time.Now().Format(time.RFC3339)
		println(dt+" "+statement)
	}
}