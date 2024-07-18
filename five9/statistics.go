package five9

import (
	"context"
	"fmt"
	"net/http"
)

type StatisticsService struct {
	authState *authenticationState
}

func (s *StatisticsService) GetRecordingbyId(ctx context.Context, agentID uint64, recordingID string) ([]byte, error) {
	var target []byte

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/strsvcs/rs/svc/agents/%d/recordings/%s?download=true", agentID, recordingID),
		http.NoBody,
	)

	if err != nil {
		return nil, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}
