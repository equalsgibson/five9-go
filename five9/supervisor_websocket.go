package five9

import (
	"context"
	"errors"
	"fmt"
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
			agentInfo:   map[five9types.UserID]five9types.AgentInfo{},
			reasonCodes: map[five9types.ReasonCodeID]five9types.ReasonCodeInfo{},
		}
	}

	websocketError := make(chan error)
	webSocketCacheReady := make(chan bool)

	{ // Ping handling
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		go func() {
			_ = s.websocketHandler.Write(ctx, []byte("ping"))
			for {
				select {
				case <-ticker.C:
					_ = s.websocketHandler.Write(ctx, []byte("ping"))
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Pong monitoring
	{
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		go func() {
			select {
			case <-ticker.C:
				if time.Since(*s.webSocketCache.lastPong) > time.Second*45 {
					websocketError <- errors.New("last valid ping response from WS is older than 45 seconds, closing connection")
				}
			case <-ctx.Done():
				return
			}
		}()
	}

	// Warm up domainMetadata cache
	go func() {
		agents, err := s.GetAllDomainUsers(ctx) // Could take a long time (6 seconds)
		if err != nil {
			websocketError <- err

			return
		}

		for _, agent := range agents {
			s.domainMetadataCache.agentInfo[agent.ID] = agent
		}

		webSocketCacheReady <- true
	}()

	// Get full statistics
	go func() {

		// Cannot use an if statement here, as we would be waiting for the context to be done (receive a message)
		select {
		case <-ctx.Done():
			return
		case <-webSocketCacheReady:
			// When starting a new session, this is called by Five9. Account for rejoining an existing session by also
			// calling this.
			if err := s.requestWebSocketFullStatistics(ctx); err != nil {
				websocketError <- err
			}
		}
	}()

	// Forever read the bytes
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
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

		}
	}()

	return <-websocketError
}

func (s *SupervisorService) AgentState(ctx context.Context) (map[five9types.UserName]five9types.AgentState, error) {
	response := map[five9types.UserName]five9types.AgentState{}

	for agentID, agentState := range s.webSocketCache.agentState {
		agentInfo, ok := s.domainMetadataCache.agentInfo[agentID]
		if !ok {
			continue
		}

		response[agentInfo.UserName] = agentState
	}

	return response, nil
}
