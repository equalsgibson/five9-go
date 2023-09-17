package five9

type SupervisorService struct {
	client    *client
	webSocket *webSocketService
}

func (s SupervisorService) WebSocket() *webSocketService {
	return s.webSocket
}
