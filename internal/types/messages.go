package types

type Message struct {
	Id             int64  `json:"id"`
	Text           string `json:"text"`
	Animal         string `json:"animal"`
	ConversationId int64  `json:"conversation_id"`
	SenderID       int64  `json:"sender_id"`
	CreatedTime    string `json:"created_time"`
	UpdatedTime    string `json:"updated_time"`
}