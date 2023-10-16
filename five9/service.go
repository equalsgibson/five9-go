package five9

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/equalsgibson/five9-go/five9/five9types"
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

	s := &Service{
		agentService: &AgentService{
			authState: &authenticationState{
				client:         c,
				apiContextPath: agentAPIContextPath,
				loginMutex:     &sync.Mutex{},
			},
			websocketHandler:    &liveWebsocketHandler{},
			domainMetadataCache: &domainMetadata{},
			webSocketCache:      &agentWebsocketCache{},
		},
		supervisorService: &SupervisorService{
			authState: &authenticationState{
				client:         c,
				apiContextPath: supervisorAPIContextPath,
				loginMutex:     &sync.Mutex{},
			},
			websocketHandler:    &liveWebsocketHandler{},
			domainMetadataCache: &domainMetadata{},
			webSocketCache:      &supervisorWebsocketCache{},
		},
	}

	for _, configFunc := range configFuncs {
		configFunc(s)
	}

	return s
}

type Service struct {
	agentService      *AgentService
	supervisorService *SupervisorService
}

func (s *Service) Supervisor() *SupervisorService {
	return s.supervisorService
}

func (s *Service) Agent() *AgentService {
	return s.agentService
}

type domainMetadata struct {
	reasonCodes    map[five9types.ReasonCodeID]five9types.ReasonCodeInfo
	agentInfoState agentInfoState
}

type agentInfoState struct {
	agentInfo   map[five9types.UserID]five9types.AgentInfo
	mutex       *sync.Mutex
	lastUpdated *time.Time
}
