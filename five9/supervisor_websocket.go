package five9

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
	"github.com/equalsgibson/five9-go/five9/internal/utils"
	"github.com/google/uuid"
)

type supervisorWebSocketCache struct {
	agentState *utils.MemoryCacheInstance[
		five9types.UserID,
		five9types.AgentState,
	]
	agentStatistics *utils.MemoryCacheInstance[
		five9types.UserID,
		five9types.AgentStatistics,
	]
	acdState *utils.MemoryCacheInstance[
		five9types.QueueID,
		five9types.ACDState,
	]
	timers *utils.MemoryCacheInstance[
		five9types.EventID,
		*time.Time,
	]
}

func (s *SupervisorService) StartWebsocket(parentCtx context.Context) error {
	// Clear any stale data from a previous connection
	s.resetCache()

	// If we encounter an error on the WebsocketErr channel, cancel the context, thus cancelling all other goroutines.
	ctx, cancel := context.WithCancel(parentCtx)

	defer func() {
		// Clear the cache when closing the connection
		s.resetCache()
		s.webSocketHandler.Close()
		cancel()
	}()

	login, err := s.authState.getLogin(ctx)
	if err != nil {
		return err
	}

	connectionURL := fmt.Sprintf("wss://%s/supsvcs/sws/%s", login.GetAPIHost(), uuid.NewString())

	if err := s.webSocketHandler.Connect(ctx, connectionURL, s.authState.client.httpClient); err != nil {
		return err
	}

	websocketError := make(chan error)

	// Ping intervals
	go func() {
		if err := s.ping(ctx); err != nil {
			websocketError <- err
		}
	}()

	// Pong monitoring
	go func() {
		if err := s.pong(ctx); err != nil {
			websocketError <- err
		}
	}()

	// Get full statistics
	go func() {
		// When starting a new session, this is called by Five9. Account for rejoining an existing session by also
		// calling this.
		if err := s.requestWebSocketFullStatistics(ctx); err != nil {
			websocketError <- err
		}
	}()

	// Forever read the bytes
	go func() {
		if err := s.read(ctx); err != nil {
			websocketError <- err
		}
	}()

	return <-websocketError
}

func (s *SupervisorService) WSAgentState(ctx context.Context) (map[five9types.UserName]five9types.AgentState, error) {
	response := map[five9types.UserName]five9types.AgentState{}

	domainUsers, err := s.getDomainUserInfoMap(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := s.webSocketCache.timers.Get(five9types.EventIDSupervisorStats); !ok {
		return nil, ErrWebSocketCacheNotReady
	}

	for agentID, agentState := range s.webSocketCache.agentState.GetAll().Items {
		agentInfo, ok := domainUsers[agentID]
		if !ok {
			continue
		}

		response[agentInfo.UserName] = agentState
	}

	return response, nil
}

func (s *SupervisorService) WSAgentStatistics(ctx context.Context) (map[five9types.UserName]five9types.AgentStatistics, error) {
	response := map[five9types.UserName]five9types.AgentStatistics{}

	_, err := s.getDomainUserInfoMap(ctx)
	if err != nil {
		return nil, err
	}

	if _, ok := s.webSocketCache.timers.Get(five9types.EventIDSupervisorStats); !ok {
		return nil, ErrWebSocketCacheNotReady
	}

	// for agentID, agentState := range s.webSocketCache.agentState.GetAll().Items {
	// 	agentInfo, ok := domainUsers[agentID]
	// 	if !ok {
	// 		continue
	// 	}

	// 	// response[agentInfo.UserName] = agentState
	// }

	return response, nil
}

func (s *SupervisorService) ping(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	if err := s.webSocketHandler.Write(ctx, []byte("ping")); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := s.webSocketHandler.Write(ctx, []byte("ping")); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *SupervisorService) pong(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lastPongReceived, ok := s.webSocketCache.timers.Get(five9types.EventIDPongReceived)
			if !ok {
				return errors.New("could not obtain last pong time from cache")
			}
			if time.Since(*lastPongReceived) > time.Second*45 {
				return errors.New("last valid ping response from WS is older than 45 seconds, closing connection")
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *SupervisorService) read(ctx context.Context) error {
	for {
		messageBytes, err := s.webSocketHandler.Read(ctx)
		if err != nil {
			return err
		}

		if err := s.handleWebsocketMessage(messageBytes); err != nil {
			return err
		}
	}
}

func (s *SupervisorService) resetCache() {
	s.webSocketCache.acdState.Reset()
	s.webSocketCache.agentState.Reset()
	s.webSocketCache.agentStatistics.Reset()
	s.webSocketCache.timers.Reset()

	s.domainMetadataCache.agentInfoState.Reset()
	s.domainMetadataCache.queueInfoState.Reset()
	s.domainMetadataCache.reasonCodeInfoState.Reset()

	serviceReset := time.Now()
	s.webSocketCache.timers.Update(five9types.EventIDPongReceived, &serviceReset)
}
