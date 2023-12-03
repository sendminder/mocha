package types

type Message struct {
	ID          int64  `json:"id"`
	Text        string `json:"text"`
	Animal      string `json:"animal"`
	ChannelID   int64  `json:"channel_id"`
	SenderID    int64  `json:"sender_id"`
	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}
