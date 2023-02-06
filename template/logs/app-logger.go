package logs

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func AppLog(message string, trackId string, threadName string, metadata interface{}) {
	log.WithFields(log.Fields{
		"applicationName": "airtime-ussd",
		"@timestamp":      time.Now(),
		"message":         message,
		"trackId":         trackId,
		"threadName":      threadName,
		"metadata":        metadata,
	}).Info(message)

}
