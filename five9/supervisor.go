package five9

import (
	"context"
	"net/http"

	"github.com/equalsgibson/five9-go/five9/five9types"
)

type SupervisorService struct {
	authState           *authenticationState
	webSocketHandler    webSocketHandler
	webSocketCache      *supervisorWebSocketCache
	domainMetadataCache *domainMetadataCache
}

func (s *SupervisorService) GetOwnUserInfo(ctx context.Context) (five9types.AgentInfo, error) {
	users, err := s.getDomainUserInfoMap(ctx)
	if err != nil {
		return five9types.AgentInfo{}, err
	}

	self, ok := users[s.authState.loginResponse.UserID]
	if !ok {
		return five9types.AgentInfo{}, ErrUnknownUserID
	}

	return self, nil
}

func (s *SupervisorService) GetStatisticsFilterSettings(ctx context.Context) ([]five9types.AgentInfo, error) {
	var target []five9types.AgentInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/supervisors/:userID/stats_filter_settings",
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

func (s *SupervisorService) SetStatisticsFilterSettings(ctx context.Context, payload any) ([]five9types.AgentInfo, error) {
	var target []five9types.AgentInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/supsvcs/rs/svc/supervisors/:userID/stats_filter_settings",
		structToReaderCloser(payload),
	)
	if err != nil {
		return nil, err
	}

	if err := s.authState.requestWithAuthentication(request, &target); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *SupervisorService) getDomainUserInfoMap(ctx context.Context) (map[five9types.UserID]five9types.AgentInfo, error) {
	allUserInfo, err := s.domainMetadataCache.agentInfoState.GetAll()
	if err == nil {
		return allUserInfo.Items, nil
	}

	domainUserSlice, err := s.GetAllDomainUsers(ctx)
	if err != nil {
		return nil, err
	}

	freshData := map[five9types.UserID]five9types.AgentInfo{}

	for _, domainUser := range domainUserSlice {
		freshData[domainUser.ID] = domainUser
	}

	if err := s.domainMetadataCache.agentInfoState.Replace(freshData); err != nil {
		return nil, err
	}

	return freshData, nil
}

func (s *SupervisorService) GetAllDomainUsers(ctx context.Context) ([]five9types.AgentInfo, error) {
	var target []five9types.AgentInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/users",
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

func (s *SupervisorService) getQueueInfoMap(ctx context.Context) (map[five9types.QueueID]five9types.QueueInfo, error) {
	q, err := s.domainMetadataCache.queueInfoState.GetAll()
	if err == nil {
		return q.Items, nil
	}

	queues, err := s.GetAllQueues(ctx)
	if err != nil {
		return nil, err
	}

	freshData := map[five9types.QueueID]five9types.QueueInfo{}

	for _, queue := range queues {
		freshData[queue.ID] = queue
	}

	if err := s.domainMetadataCache.queueInfoState.Replace(freshData); err != nil {
		return nil, err
	}

	return freshData, nil
}

func (s *SupervisorService) GetAllQueues(ctx context.Context) ([]five9types.QueueInfo, error) {
	var target []five9types.QueueInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/skills",
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

func (s *SupervisorService) GetReasonCodeInfoMap(ctx context.Context) (map[five9types.ReasonCodeID]five9types.ReasonCodeInfo, error) {
	r, err := s.domainMetadataCache.reasonCodeInfoState.GetAll()
	if err == nil {
		return r.Items, nil
	}

	reasonCodes, err := s.GetAllReasonCodes(ctx)
	if err != nil {
		return nil, err
	}

	freshData := map[five9types.ReasonCodeID]five9types.ReasonCodeInfo{}

	for _, reasonCode := range reasonCodes {
		freshData[reasonCode.ID] = reasonCode
	}

	if err := s.domainMetadataCache.reasonCodeInfoState.Replace(freshData); err != nil {
		return nil, err
	}

	return freshData, nil
}

func (s *SupervisorService) GetAllReasonCodes(ctx context.Context) ([]five9types.ReasonCodeInfo, error) {
	reasonCodes := []five9types.ReasonCodeInfo{}

	logoutCodes, err := s.getAllLogoutReasonCodes(ctx)
	if err != nil {
		return nil, err
	}

	reasonCodes = append(reasonCodes, logoutCodes...)

	notReadyCodes, err := s.getAllNotReadyReasonCodes(ctx)
	if err != nil {
		return nil, err
	}

	return append(reasonCodes, notReadyCodes...), nil
}

func (s *SupervisorService) getAllLogoutReasonCodes(ctx context.Context) ([]five9types.ReasonCodeInfo, error) {
	var target []five9types.ReasonCodeInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/logout_reason_codes",
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

func (s *SupervisorService) getAllNotReadyReasonCodes(ctx context.Context) ([]five9types.ReasonCodeInfo, error) {
	var target []five9types.ReasonCodeInfo

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/supsvcs/rs/svc/orgs/:organizationID/not_ready_reason_codes",
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

func (s *SupervisorService) requestWebSocketFullStatistics(ctx context.Context) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"/supsvcs/rs/svc/supervisors/:userID/request_full_statistics",
		http.NoBody,
	)
	if err != nil {
		return err
	}

	if err := s.authState.requestWithAuthentication(request, nil); err != nil {
		return err
	}

	return nil
}

func (s *SupervisorService) UpdateAgentState(ctx context.Context, agentID five9types.UserID) (five9types.UserFullStateInfo, error) {
	return five9types.UserFullStateInfo{}, nil
}
