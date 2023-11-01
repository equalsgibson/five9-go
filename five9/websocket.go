package five9

import (
	"context"
	"errors"
	"math"
	"net/http"

	"nhooyr.io/websocket"
)

type webSocketHandler interface {
	Connect(ctx context.Context, s string, c *http.Client) error
	Read(ctx context.Context) ([]byte, error)
	Write(ctx context.Context, b []byte) error
	Close()
}

type liveWebsocketHandler struct {
	c *websocket.Conn
}

func (h *liveWebsocketHandler) Connect(ctx context.Context, connectionURL string, httpClient *http.Client) error {
	if h.c != nil {
		h.c.Close(websocket.StatusNormalClosure, "Restarting connection")
	}

	h.c = nil

	connection, response, err := websocket.Dial(ctx, connectionURL, &websocket.DialOptions{
		HTTPClient: httpClient,
	})
	if err != nil {
		return err
	}

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	// Set a very large read limit
	connection.SetReadLimit(math.MaxInt32)

	h.c = connection

	return nil
}

func (h *liveWebsocketHandler) Read(ctx context.Context) ([]byte, error) {
	messageType, b, err := h.c.Read(ctx)
	if err != nil {
		return nil, err
	}

	if messageType != websocket.MessageText {
		return nil, errors.New("binary messages are not supported")
	}

	return b, nil
}

func (h *liveWebsocketHandler) Close() {
	if h.c != nil {
		h.c.Close(websocket.StatusNormalClosure, "closed")
	}
}

func (h *liveWebsocketHandler) Write(ctx context.Context, data []byte) error {
	return h.c.Write(ctx, websocket.MessageText, data)
}
