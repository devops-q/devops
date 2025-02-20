package json_models

import "time"

type Message struct {
	Content string    `json:"content"`
	PubDate time.Time `json:"pub_date"`
	User    string    `json:"user"`
}

type CreateMessageBody struct {
	Content string `json:"content"`
}
