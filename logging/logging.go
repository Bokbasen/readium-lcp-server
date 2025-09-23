// Copyright 2023 European Digital Reading Lab. All rights reserved.
// Licensed to the Readium Foundation under one or more contributor license agreements.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file exposed on Github (readium) in the project repository.

package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/readium/readium-lcp-server/config"
	"github.com/slack-go/slack"
)

var (
	LogFile        *log.Logger
	SlackApi       *slack.Client
	SlackChannelID string
)

// Init inits the log file and opens it
func Init(logging config.Logging) error {
	//logPath string, cm bool
	if logging.Directory != "" {
		log.Println("Open log file as " + logging.Directory)
		file, err := os.OpenFile(logging.Directory, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		LogFile = log.New(file, "", log.LUTC)
	}
	if logging.SlackToken != "" && logging.SlackChannelID != "" {
		log.Println("Init Slack connection")
		SlackApi = slack.New(logging.SlackToken)
		if SlackApi == nil {
			return errors.New("error creating a Slack connector")
		}
		SlackChannelID = logging.SlackChannelID
	}
	return nil
}

func Print(message string)                   { Info(message) }
func Println(v ...interface{})               { Infoln(v...) }
func Printf(format string, v ...interface{}) { Infof(format, v...) }

func writeLog(level string, message string) {
	// log on stdout
	log.Printf("[%s] %s", level, message)

	// log on a file
	if LogFile != nil {
		currentTime := time.Now().UTC().Format(time.RFC3339)
		LogFile.Println(currentTime + "\t[" + level + "]\t" + message)
	}

	// log on Slack (only for WARN and ERROR to reduce noise)
	if SlackApi != nil && (level == "WARN" || level == "ERROR") {
		_, _, err := SlackApi.PostMessage(SlackChannelID, slack.MsgOptionText("["+level+"] "+message, false))
		if err != nil {
			log.Printf("Error sending Slack msg, %v", err)
		}
	}
}

func Debug(message string) { writeLog("DEBUG", message) }
func Info(message string)  { writeLog("INFO", message) }
func Warn(message string)  { writeLog("WARN", message) }
func Error(message string) { writeLog("ERROR", message) }

// Formatted
func Debugf(format string, v ...interface{}) { writeLog("DEBUG", fmt.Sprintf(format, v...)) }
func Infof(format string, v ...interface{})  { writeLog("INFO", fmt.Sprintf(format, v...)) }
func Warnf(format string, v ...interface{})  { writeLog("WARN", fmt.Sprintf(format, v...)) }
func Errorf(format string, v ...interface{}) { writeLog("ERROR", fmt.Sprintf(format, v...)) }

// Println
func Debugln(v ...interface{}) { writeLog("DEBUG", fmt.Sprintln(v...)) }
func Infoln(v ...interface{})  { writeLog("INFO", fmt.Sprintln(v...)) }
func Warnln(v ...interface{})  { writeLog("WARN", fmt.Sprintln(v...)) }
func Errorln(v ...interface{}) { writeLog("ERROR", fmt.Sprintln(v...)) }
