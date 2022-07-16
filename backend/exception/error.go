package exception

import log "github.com/sirupsen/logrus"

func Check(err error) {
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
}

type MessageModel struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(message string) MessageModel {
	panic(MessageModel{Message: message})
}
