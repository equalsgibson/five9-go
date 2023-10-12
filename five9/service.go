package five9

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
)

func NewService(
	creds PasswordCredentials,
	configFuncs ...ConfigFunc,
) *Service {
	cookieJar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar: cookieJar,
	}

	c := &client{
		credentials: creds,
		httpClient:  httpClient,
		agentAuth: &authenticationState{
			apiContextPath: agentAPIContextPath,
			loginMutex:     &sync.Mutex{},
		},
		supervisorAuth: &authenticationState{
			apiContextPath: supervisorAPIContextPath,
			loginMutex:     &sync.Mutex{},
		},
		requestPreProcessors: []func(r *http.Request) error{},
	}

	s := &Service{
		client: c,
		agentService: &AgentService{
			client:              c,
			websocketHandler:    &liveWebsocketHandler{},
			websocketReady:      make(chan bool),
			domainMetadataCache: &domainMetadata{},
			webSocketCache:      &agentWebsocketCache{},
		},
		supervisorService: &SupervisorService{
			client:              c,
			websocketHandler:    &liveWebsocketHandler{},
			websocketReady:      make(chan bool),
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
	client            *client
	agentService      *AgentService
	supervisorService *SupervisorService
}

func (s *Service) Supervisor() *SupervisorService {
	return s.supervisorService
}

func (s *Service) Agent() *AgentService {
	return s.agentService
}
