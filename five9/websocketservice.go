package five9

import "golang.org/x/net/websocket"

type webSocketService struct {
	client *client
}

func (w *webSocketService) Connect() (*websocket.Conn, error) {
	connection, err := websocket.Dial("asd", "asd", "asd")
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	return connection, nil
}

func (w *webSocketService) Ping() {
}

func (w *webSocketService) Close() {
}

// func (w *webSocketService) reader(conn *websocket.Conn) error {
// 	// Listen indefinitely
// 	// for {
// 	// 	messageType, p, err := conn.
// 	// }
// 	return nil
// }
