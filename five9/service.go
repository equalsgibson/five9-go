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
		credentials:          creds,
		httpClient:           httpClient,
		loginMutex:           &sync.Mutex{},
		requestPreProcessors: []func(r *http.Request) error{},
	}

	s := &Service{
		client: c,
		supervisorService: &SupervisorService{
			client:           c,
			websocketHandler: &liveWebsocketHandler{},
			websocketReady:   make(chan bool),
		},
	}

	for _, configFunc := range configFuncs {
		configFunc(s)
	}

	return s
}

type Service struct {
	client            *client
	supervisorService *SupervisorService
}

func (s *Service) Supervisor() *SupervisorService {
	return s.supervisorService
}
