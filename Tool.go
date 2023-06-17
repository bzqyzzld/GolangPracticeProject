package main

import (
	"fmt"
	"time"
)

func DealMsg(msg string, u interface{}) string {
	user, ok := u.(*User)
	var name string
	if ok == true {
		name = user.UserName
	} else {
		name = "system"
	}
	nowString := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s[%s]%s", nowString, name, msg)
}
