package models

type MessageModel struct {
    Number   string          `json:"number"`
    Message  string          `json:"message"`
    SenderId string          `json:"senderid"`
    Priority MessagePriority `json:"priority"`
}
