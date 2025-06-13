package five9

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StatisticsService struct {
	authState *authenticationState
}

func (s *StatisticsService) GetRecordingbyId(ctx context.Context, agentID uint64, recordingID string) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/strsvcs/rs/svc/agents/%d/recordings/%s?download=true", agentID, recordingID),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	response, err := s.authState.requestDownloadWithAuthentication(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	target, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return target, nil
}
