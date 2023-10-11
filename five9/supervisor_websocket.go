package five9

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type supervisorWebsocketCache struct {
	agentState map[UserID]AgentState
	lastPong   *time.Time
}

func (s *SupervisorService) StartWebsocket(ctx context.Context) error {
	login, err := s.client.getLogin(ctx)
	if err != nil {
		return err
	}

	connectionURL := fmt.Sprintf("wss://%s/supsvcs/sws/%s", login.GetAPIHost(), uuid.NewString())

	if err := s.websocketHandler.Connect(ctx, connectionURL, s.client.httpClient); err != nil {
		return err
	}

	defer func() {
		s.websocketHandler.Close()
	}()

	{ // reset state upon starting the websocket connection
		now := time.Now()
		s.websocketReady = make(chan bool)
		s.webSocketCache = &supervisorWebsocketCache{
			agentState: map[UserID]AgentState{},
			lastPong:   &now,
		}
		s.domainMetadataCache = &domainMetadata{
			agentInfo:   map[UserID]AgentInfo{},
			reasonCodes: map[ReasonCodeID]ReasonCodeInfo{},
		}
	}

	websocketError := make(chan error)

	{ // Ping handling
		ticker := time.NewTicker(time.Second * 5)
		go func() {
			_ = s.websocketHandler.Write(ctx, []byte("ping"))
			for range ticker.C {
				_ = s.websocketHandler.Write(ctx, []byte("ping"))
			}
		}()

		defer ticker.Stop()
	}

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

	// Getting all users (only a go routine because this call can be slow)
	go func() {
		agents, err := s.getAllDomainUsers(ctx) // Could take a long time (6 seconds)
		if err != nil {
			websocketError <- err

			return
		}

		for _, agent := range agents {
			s.domainMetadataCache.agentInfo[agent.ID] = agent
		}
	}()

	// Get Domain Metadata

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

func (s *SupervisorService) AgentState(ctx context.Context) (map[UserName]AgentState, error) {
	response := map[UserName]AgentState{}

	for agentID, agentState := range s.webSocketCache.agentState {
		agentInfo, ok := s.domainMetadataCache.agentInfo[agentID]
		if !ok {
			continue
		}

		response[agentInfo.UserName] = agentState
	}

	return response, nil
}
