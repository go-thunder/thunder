package events

const EmailTopic = "topic.email"

type EmailPayload struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
