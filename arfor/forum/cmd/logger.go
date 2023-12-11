package main

import (
	"log"
	"os"
)

func errLogger(err error) *log.Logger {

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return errorLog
}

func infoLog(n string) *log.Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	return infoLog

}
