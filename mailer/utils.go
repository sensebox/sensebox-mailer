package mailer

import (
	"fmt"
	"time"
)

func LogInfo(msgid string, msgs ...interface{}) {
	logPrefix := time.Now().UTC().Format(time.RFC3339) + " [" + msgid + "]"
	msgs = append([]interface{}{logPrefix}, msgs...)
	fmt.Println(msgs...)
}
