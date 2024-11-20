package events

const TestTopic = "topic.test"

type TestPayload struct {
	Hello string `json:"hello"`
}
