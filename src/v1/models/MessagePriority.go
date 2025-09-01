package models

type MessagePriority string

const (
    HIGHEST  MessagePriority = "0"
    HIGH     MessagePriority = "1"
    MEDIUM   MessagePriority = "2"
    LOW      MessagePriority = "3"
    LOWEST   MessagePriority = "4"
)
