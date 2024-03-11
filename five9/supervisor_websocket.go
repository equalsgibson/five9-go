package five9

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/equalsgibson/concur/concur"
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
	ctx, cancel := context.WithCancelCause(parentCtx)

	defer func() {
		// Clear the cache when closing the connection
		s.resetCache()
	}()

	login, err := s.authState.getLogin(ctx)
	if err != nil {
		return err
	}

	connectionURL := fmt.Sprintf("wss://%s/supsvcs/sws/%s", login.GetAPIHost(), uuid.NewString())

	if err := s.webSocketHandler.Connect(ctx, connectionURL, s.authState.client.httpClient); err != nil {
		return err
	}
	defer s.webSocketHandler.Close()

	asyncReader := concur.NewAsyncReader(s.webSocketHandler.Read)
	go asyncReader.Loop(ctx)
	defer asyncReader.Close()

	pingTicker := time.NewTicker(time.Second * 5)
	defer pingTicker.Stop()
	go func() {
		for {
			select {
			case <-pingTicker.C:
				if err := s.ping(ctx); err != nil {
					cancel(err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	pongMonitorTicker := time.NewTicker(time.Second * 5)
	defer pongMonitorTicker.Stop()
	go func() {
		for {
			select {
			case <-pingTicker.C:
				if err := s.pong(ctx); err != nil {
					cancel(err)
					return
					// asyncReader.Close()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Get full statistics
	go func() {
		// When starting a new session, this is called by Five9. Account for rejoining an existing session by also
		// calling this.
		if err := s.requestWebSocketFullStatistics(ctx); err != nil {
			cancel(err)
			// asyncReader.Close()
		}
	}()

	for {
		select {
		case update := <-asyncReader.Updates():
			if update.Err != nil {
				return update.Err
			}

			if err := s.handleWebsocketMessage(update.Item); err != nil {
				return err
			}
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}

func (s *SupervisorService) WSAgentState(ctx context.Context) (map[five9types.UserName]five9types.AgentState, error) {
	response := map[five9types.UserName]five9types.AgentState{}

	domainUsers, err := s.getDomainUserInfoMap(ctx)
	if err != nil {
		return nil, err
	}

	all, err := s.webSocketCache.agentState.GetAll()
	if err != nil {
		if errors.Is(err, utils.ErrWebSocketCacheStale) {
			return nil, ErrWebSocketCacheStale
		}

		if errors.Is(err, utils.ErrWebSocketCacheNotReady) {
			return nil, ErrWebSocketCacheNotReady
		}

		return nil, err
	}

	for agentID, agentState := range all.Items {
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

	domainUsers, err := s.getDomainUserInfoMap(ctx)
	if err != nil {
		return nil, err
	}

	allDomainUsers, err := s.webSocketCache.agentStatistics.GetAll()
	if err != nil {
		if errors.Is(err, utils.ErrWebSocketCacheStale) {
			return nil, ErrWebSocketCacheStale
		}

		if errors.Is(err, utils.ErrWebSocketCacheNotReady) {
			return nil, ErrWebSocketCacheNotReady
		}

		return nil, err
	}

	for agentID, agentStatistic := range allDomainUsers.Items {
		agentInfo, ok := domainUsers[agentID]
		if !ok {
			continue
		}

		response[agentInfo.UserName] = agentStatistic
	}

	return response, nil
}

func (s *SupervisorService) WSACDState(ctx context.Context) (map[string]five9types.ACDState, error) {
	response := map[string]five9types.ACDState{}

	queues, err := s.getQueueInfoMap(ctx)
	if err != nil {
		return nil, err
	}

	allACDState, err := s.webSocketCache.acdState.GetAll()
	if err != nil {
		if errors.Is(err, utils.ErrWebSocketCacheStale) {
			return nil, ErrWebSocketCacheStale
		}

		if errors.Is(err, utils.ErrWebSocketCacheNotReady) {
			return nil, ErrWebSocketCacheNotReady
		}

		return nil, err
	}

	for queueID, queueState := range allACDState.Items {
		queueInfo, ok := queues[queueID]
		if !ok {
			continue
		}

		response[queueInfo.Name] = queueState
	}

	return response, nil
}

func (s *SupervisorService) ping(ctx context.Context) error {
	if err := s.webSocketHandler.Write(ctx, []byte("ping")); err != nil {
		return err
	}

	return nil
}

func (s *SupervisorService) pong(_ context.Context) error {
	lastPongReceived, ok := s.webSocketCache.timers.Get(five9types.EventIDPongReceived)
	if !ok {
		return errors.New("could not obtain last pong time from cache")
	}

	if time.Since(*lastPongReceived) > time.Second*45 {
		return errors.New("last valid ping response from WS is older than 45 seconds, closing connection")
	}

	return nil
}

func (s *SupervisorService) resetCache() {
	s.authState.loginMutex.Lock()
	s.authState.loginResponse = nil
	s.authState.loginMutex.Unlock()

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
