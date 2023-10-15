package five9

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"

	"github.com/google/uuid"
)

type supervisorWebsocketCache struct {
	agentState map[five9types.UserID]five9types.AgentState
	lastPong   *time.Time
}

func (s *SupervisorService) StartWebsocket(ctx context.Context) error {
	login, err := s.authState.getLogin(ctx)
	if err != nil {
		return err
	}

	connectionURL := fmt.Sprintf("wss://%s/supsvcs/sws/%s", login.GetAPIHost(), uuid.NewString())

	if err := s.websocketHandler.Connect(ctx, connectionURL, s.authState.client.httpClient); err != nil {
		return err
	}

	defer func() {
		s.websocketHandler.Close()
	}()

	{ // reset state upon starting the websocket connection
		now := time.Now()
		s.websocketReady = make(chan bool)
		s.webSocketCache = &supervisorWebsocketCache{
			agentState: map[five9types.UserID]five9types.AgentState{},
			lastPong:   &now,
		}
		s.domainMetadataCache = &domainMetadata{
			agentInfoState: agentInfoState{
				mutex:     &sync.Mutex{},
				agentInfo: map[five9types.UserID]five9types.AgentInfo{},
			},
			reasonCodes: map[five9types.ReasonCodeID]five9types.ReasonCodeInfo{},
		}
	}

	websocketError := make(chan error)

	go s.ping(ctx, websocketError)

	// Pong monitoring
	{
		ticker := time.NewTicker(time.Second)
		go func() {
			for range ticker.C {
				if time.Since(*s.webSocketCache.lastPong) > time.Second*45 {
					websocketError <- errors.New("last valid ping response from WS is older than 45 seconds, closing connection")

					return
				}
			}
		}()

		defer ticker.Stop()
	}

	// Get full statistics
	go func() {
		<-s.websocketReady // Block until message received from channel, I don't need the value

		// When starting a new session, this is called by Five9. Account for rejoining an existing session by also
		// calling this.
		if err := s.requestWebSocketFullStatistics(ctx); err != nil {
			websocketError <- err
		}
	}()

	// Forever read the bytes
	go func() {
		for {
			messageBytes, err := s.websocketHandler.Read(ctx)
			if err != nil {
				websocketError <- err

				return
			}

			if err := s.handleWebsocketMessage(messageBytes); err != nil {
				websocketError <- err

				return
			}
		}
	}()

	return <-websocketError
}

func (s *SupervisorService) WSAgentState(ctx context.Context) (map[five9types.UserName]five9types.AgentState, error) {
	response := map[five9types.UserName]five9types.AgentState{}

	if len(s.domainMetadataCache.agentInfoState.agentInfo) < 1 {
		log.Println("domainMetadata has not been fetched, getting data")
		_, err := s.GetAllDomainUsers(ctx)
		if err != nil {
			return nil, err
		}
	}

	for agentID, agentState := range s.webSocketCache.agentState {
		agentInfo, ok := s.domainMetadataCache.agentInfoState.agentInfo[agentID]
		if !ok {
			continue
		}

		response[agentInfo.UserName] = agentState
	}

	return response, nil
}

func (s *SupervisorService) ping(ctx context.Context, errChan chan<- error) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	if err := s.websocketHandler.Write(ctx, []byte("ping")); err != nil {
		errChan <- err
	}

	for {
		select {
		case <-ticker.C:
			if err := s.websocketHandler.Write(ctx, []byte("ping")); err != nil {
				errChan <- err
			}
			log.Print("pinged!")
		case <-ctx.Done():
			log.Print("context cancelled, ending process")
			return
		}
	}
}

// { // Ping handling
// 	ticker := time.NewTicker(time.Second * 5)
// 	go func() {
// 		_ = s.websocketHandler.Write(ctx, []byte("ping"))
// 		for range ticker.C {
// 			_ = s.websocketHandler.Write(ctx, []byte("ping"))
// 		}
// 	}()

// 	defer ticker.Stop()
// }
