package exception

import "log"

func Check(err error) {
	if err != nil {
		log.Println(err)
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
