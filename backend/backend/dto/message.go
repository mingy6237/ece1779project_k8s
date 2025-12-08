package dto

type Message struct {
	SenderInstanceID string      `json:"sender_instance_id"`
	Payload          interface{} `json:"payload"`
}
