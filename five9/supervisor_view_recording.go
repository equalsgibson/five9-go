package five9

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RecordingRequestPayload struct {
	Limit        uint   `json:"limit"`
	SortField    string `json:"sortField"`
	Ascending    bool   `json:"ascending"`
	ShowUploaded bool   `json:"showUploaded"`
}

type Record struct {
	ID            string `json:"id"`
	CampaignID    string `json:"campaignId"`
	Created       int64  `json:"created"`
	Number        string `json:"number"`
	Name          string `json:"name"`
	Length        int64  `json:"length"`
	Status        string `json:"status"`
	CallSessionID string `json:"callSessionId"`
}

type RecordsResponse struct {
	ID      string   `json:"id"`
	Records []Record `json:"records"`
}

func (s *SupervisorService) GetRecordingId(ctx context.Context, agentID uint64) ([]Record, error) {
	var target RecordsResponse

	requestBody := RecordingRequestPayload{
		Limit:        100,
		SortField:    "CREATED",
		Ascending:    false,
		ShowUploaded: true,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("/supsvcs/rs/svc/supervisors/:userID/agents/%d/recording_views", agentID),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target.Records, nil
}
