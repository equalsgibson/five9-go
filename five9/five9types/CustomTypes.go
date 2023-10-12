package five9types

type WebsocketMessage struct {
	Context WebsocketMessageContext `json:"context"`
	Payload any                     `json:"payLoad"`
}

type WebsocketMessageContext struct {
	EventID EventID `json:"eventId"`
}
