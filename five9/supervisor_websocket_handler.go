package five9

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type websocketFrameProcessingError struct {
	OriginalError error
	MessageBytes  []byte
}

func (err websocketFrameProcessingError) Error() string {
	return fmt.Sprintf("Error while processing websocket frame: %s - %s", err.OriginalError.Error(), string(err.MessageBytes))
}

func (s *SupervisorService) handleWebsocketMessage(messageBytes []byte) error {
	message := five9types.WebsocketMessage{}
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

	eventReceivedTime := time.Now()
	s.webSocketCache.timers.Update(message.Context.EventID, &eventReceivedTime)

	switch message.Context.EventID {
	case five9types.EventIDServerConnected:
		return nil
	case five9types.EventIDPongReceived:
		return s.handlerPong(message.Payload)
	case five9types.EventIDIncrementalStatsUpdate:
		return s.handlerIncrementalStatsUpdate(message.Payload)
	case five9types.EventIDSupervisorStats:
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

		dataSource := five9types.DataSource(dataSourceString)

		payloadItemBytes, err := json.Marshal(payloadItem)
		if err != nil {
			return err
		}

		if dataSource == five9types.DataSourceAgentState {
			eventTarget := five9types.WebSocketIncrementalAgentStateData{}
			if err := json.Unmarshal(payloadItemBytes, &eventTarget); err != nil {
				return websocketFrameProcessingError{
					OriginalError: err,
					MessageBytes:  payloadItemBytes,
				}
			}

			if err := s.handleAgentStateUpdate(eventTarget); err != nil {
				return err
			}
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

		dataSource := five9types.DataSource(dataSourceString)

		payloadItemBytes, err := json.Marshal(payloadItem)
		if err != nil {
			return err
		}

		switch dataSource {
		// ** //
		case five9types.DataSourceAgentState:
			eventTarget := five9types.WebsocketSupervisorStateData{}
			if err := json.Unmarshal(payloadItemBytes, &eventTarget); err != nil {
				return websocketFrameProcessingError{
					OriginalError: err,
					MessageBytes:  payloadItemBytes,
				}
			}

			freshData := map[five9types.UserID]five9types.AgentState{}

			for _, agent := range eventTarget.Data {
				freshData[agent.ID] = agent
			}

			s.webSocketCache.agentState.Replace(freshData)
		// ** //
		case five9types.DataSourceAgentStatistic:
			eventTarget := five9types.WebsocketSupervisorStatisticsData{}
			if err := json.Unmarshal(payloadItemBytes, &eventTarget); err != nil {
				return websocketFrameProcessingError{
					OriginalError: err,
					MessageBytes:  payloadItemBytes,
				}
			}

			freshData := map[five9types.UserID]five9types.AgentStatistics{}

			for _, agentStatistic := range eventTarget.Data {
				freshData[agentStatistic.ID] = agentStatistic
			}

			s.webSocketCache.agentStatistics.Replace(freshData)
		// ** //
		case five9types.DataSourceACDStatus:
			eventTarget := five9types.WebsocketSupervisorACDData{}
			if err := json.Unmarshal(payloadItemBytes, &eventTarget); err != nil {
				return websocketFrameProcessingError{
					OriginalError: err,
					MessageBytes:  payloadItemBytes,
				}
			}

			freshData := map[five9types.QueueID]five9types.ACDState{}
			for _, acd := range eventTarget.Data {
				freshData[acd.ID] = acd
			}

			s.webSocketCache.acdState.Replace(freshData)
		}
	}

	return nil
}

func (s *SupervisorService) handleAgentStateUpdate(eventData five9types.WebSocketIncrementalAgentStateData) error {
	for _, addedData := range eventData.Added {
		s.webSocketCache.agentState.Update(addedData.ID, addedData)
	}

	for _, updatedData := range eventData.Updated {
		s.webSocketCache.agentState.Update(updatedData.ID, updatedData)
	}

	for _, removedData := range eventData.Removed {
		s.webSocketCache.agentState.Delete(removedData.ID)
	}

	return nil
}
