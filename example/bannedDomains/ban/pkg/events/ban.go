package events

const BanTopic = "topic.ban"

type BanPayload struct {
	ID int `json:"id"`
}
