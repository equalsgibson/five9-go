package five9

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type websocketFrameProcessingError struct {
	OriginalError error
	MessageBytes  []byte
}

func (err websocketFrameProcessingError) Error() string {
	return fmt.Sprintf("Error while processing websocket frame: %s - %s", err.OriginalError.Error(), string(err.MessageBytes))
}

type websocketMessage struct {
	Context struct {
		EventID eventID `json:"eventId"`
	} `json:"context"`
	Payload any `json:"payload"`
}

func (s *SupervisorService) handleWebsocketMessage(messageBytes []byte) error {
	message := websocketMessage{}
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return websocketFrameProcessingError{
			OriginalError: err,
			MessageBytes:  messageBytes,
		}
	}

	if message.Context.EventID == "" {
		return websocketFrameProcessingError{
			OriginalError: errors.New("unsupported message"),
			MessageBytes:  messageBytes,
		}
	}

	switch message.Context.EventID {
	case eventIDServerConnected:
		s.websocketReady <- true

		return nil
	case eventIDPongReceived:
		return s.handlerPong(message.Payload)
	case eventIDIncrementalStatsUpdate:
		return s.handlerIncrementalStatsUpdate(message.Payload)
	case eventIDSupervisorStats:
		return s.handlerSupervisorStats(message.Payload)
	}

	return nil
}

func (s *SupervisorService) handlerPong(payload any) error {
	payloadString, ok := payload.(string)
	if !ok {
		return fmt.Errorf("failed type assertion for payload: %T", payload)
	}

	if payloadString != "pong" {
		return fmt.Errorf("payload not expected type")
	}

	now := time.Now()
	s.cache.lastPong = &now

	return nil
}

func (s *SupervisorService) handlerIncrementalStatsUpdate(payload any) error {
	payloadSlice, ok := payload.([]any)
	if !ok {
		return fmt.Errorf("failed type assertion for payload: %T", payload)
	}

	for _, payloadItem := range payloadSlice {
		payloadMap, ok := payloadItem.(map[string]any)
		if !ok {
			return fmt.Errorf("failed type assertion for payload item: %T", payloadItem)
		}

		dataSourceRaw, ok := payloadMap["dataSource"]
		if !ok {
			return fmt.Errorf("data source not found: %+v", payloadItem)
		}

		dataSourceString, ok := dataSourceRaw.(string)
		if !ok {
			return fmt.Errorf("data source is not a string, %T", dataSourceRaw)
		}

		dataSource := dataSource(dataSourceString)

		payloadItemBytes, err := json.Marshal(payloadItem)
		if err != nil {
			return err
		}

		switch dataSource {
		case dataSourceAgentState:
			eventTarget := webSocketIncrementalStatsUpdateData{}
			if err := json.Unmarshal(payloadItemBytes, &eventTarget); err != nil {
				return websocketFrameProcessingError{
					OriginalError: err,
					MessageBytes:  payloadItemBytes,
				}
			}

			if err := s.handleAgentStateUpdate(eventTarget); err != nil {
				return err
			}

		case dataSourceACDStatus:
		default:
			return errors.New("unsupported")
		}

	}

	return nil
}

func (s *SupervisorService) handlerSupervisorStats(payload any) error {
	payloadSlice, ok := payload.([]any)
	if !ok {
		return fmt.Errorf("failed type assertion for payload: %T", payload)
	}

	for _, payloadItem := range payloadSlice {
		payloadMap, ok := payloadItem.(map[string]any)
		if !ok {
			return fmt.Errorf("failed type assertion for payload item: %T", payloadItem)
		}

		dataSourceRaw, ok := payloadMap["dataSource"]
		if !ok {
			return fmt.Errorf("data source not found: %+v", payloadItem)
		}

		dataSourceString, ok := dataSourceRaw.(string)
		if !ok {
			return fmt.Errorf("data source is not a string, %T", dataSourceRaw)
		}

		dataSource := dataSource(dataSourceString)

		payloadItemBytes, err := json.Marshal(payloadItem)
		if err != nil {
			return err
		}

		if dataSource == dataSourceAgentState {
			eventTarget := websocketSupervisorStatsData{}
			if err := json.Unmarshal(payloadItemBytes, &eventTarget); err != nil {
				return websocketFrameProcessingError{
					OriginalError: err,
					MessageBytes:  payloadItemBytes,
				}
			}

			for _, agent := range eventTarget.Data {
				s.cache.agentState[agent.ID] = agent
			}

			continue
		}
	}

	return nil
}

func (s *SupervisorService) handleAgentStateUpdate(eventData webSocketIncrementalStatsUpdateData) error {
	for _, addedData := range eventData.Added {
		// TODO: confirm what data looks like when agents are added
		s.cache.agentState[addedData.ID] = addedData
	}

	for _, updatedData := range eventData.Updated {
		// TODO: confirm what data looks like when agents are updated
		s.cache.agentState[updatedData.ID] = updatedData
	}

	for _, removedData := range eventData.Removed {
		// TODO: verify what the payload looks like when users are removed
		delete(s.cache.agentState, removedData.ID)
	}

	return nil
}
