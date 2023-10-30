package five9

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
	"github.com/equalsgibson/five9-go/five9/internal/utils"
)

func NewService(
	creds five9types.PasswordCredentials,
	configFuncs ...ConfigFunc,
) *Service {
	cookieJar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar: cookieJar,
	}

	c := &client{
		credentials:          creds,
		httpClient:           httpClient,
		requestPreProcessors: []func(r *http.Request) error{},
	}

	defaultLoginURL := "app.five9.com"

	s := &Service{
		// ** //
		agentService: &AgentService{
			authState: &authenticationState{
				client:         c,
				apiContextPath: agentAPIContextPath,
				loginMutex:     &sync.Mutex{},
				loginURL:       defaultLoginURL,
			},
		},
		// ** //
		supervisorService: &SupervisorService{
			authState: &authenticationState{
				client:         c,
				apiContextPath: supervisorAPIContextPath,
				loginMutex:     &sync.Mutex{},
				loginURL:       defaultLoginURL,
			},
			domainMetadataCache: &domainMetadataCache{
				agentInfoState: utils.NewMemoryCacheInstance[
					five9types.UserID,
					five9types.AgentInfo,
				](),
				reasonCodeInfoState: utils.NewMemoryCacheInstance[
					five9types.ReasonCodeID,
					five9types.ReasonCodeInfo,
				](),
				queueInfoState: utils.NewMemoryCacheInstance[
					five9types.QueueID,
					five9types.QueueInfo,
				](),
			},
			webSocketHandler: &liveWebsocketHandler{},
			webSocketCache: &supervisorWebSocketCache{
				agentState: utils.NewMemoryCacheInstance[
					five9types.UserID,
					five9types.AgentState,
				](),
				agentStatistics: utils.NewMemoryCacheInstance[
					five9types.UserID,
					five9types.AgentStatistics,
				](),
				acdState: utils.NewMemoryCacheInstance[
					five9types.QueueID,
					five9types.ACDState,
				](),
				timers: utils.NewMemoryCacheInstance[
					five9types.EventID,
					*time.Time,
				](),
			},
		},
	}

	// Set the cache to default values
	s.supervisorService.resetCache()

	for _, configFunc := range configFuncs {
		configFunc(s)
	}

	return s
}

type Service struct {
	agentService      *AgentService
	supervisorService *SupervisorService
	loginURL          string
}

func (s *Service) Supervisor() *SupervisorService {
	return s.supervisorService
}

func (s *Service) Agent() *AgentService {
	return s.agentService
}

type domainMetadataCache struct {
	reasonCodeInfoState *utils.MemoryCacheInstance[five9types.ReasonCodeID, five9types.ReasonCodeInfo]
	agentInfoState      *utils.MemoryCacheInstance[five9types.UserID, five9types.AgentInfo]
	queueInfoState      *utils.MemoryCacheInstance[five9types.QueueID, five9types.QueueInfo]
}
