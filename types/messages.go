package types

type Message struct {
	Id             int64  `json:"id"`
	Text           string `json:"text"`
	Animal         string `json:"animal"`
	Encrypted      bool   `json:"encrypted"`
	ConversationId int64  `json:"conversation_id"`
	SenderID       int64  `json:"sender_id"`
	CreatedTime    string // Modified to be string type
	UpdatedTime    string // Modified to be string type
}
