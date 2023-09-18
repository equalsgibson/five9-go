package five9

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

func NewService(
	username string,
	password string,
) *Service {
	// Cookie Jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	supervisorClient := &client{
		httpClient: &http.Client{
			Jar: jar,
		},
		authenticationMutex: &sync.Mutex{},
		credentials: PasswordCredentials{
			Username: username,
			Password: password,
		},
		context: "supsvcs/rs/svc",
	}

	agentClient := &client{
		httpClient: &http.Client{
			Jar: jar,
		},
		authenticationMutex: &sync.Mutex{},
		credentials: PasswordCredentials{
			Username: username,
			Password: password,
		},
		context: "appsvcs/rs/svc",
	}

	if err := supervisorClient.authenticate(context.Background()); err != nil {
		panic(err)
	}

	return &Service{
		agentService: &AgentService{
			client: agentClient,
			webSocket: &webSocketService{
				client: agentClient,
			},
		},
		supervisorService: &SupervisorService{
			client: supervisorClient,
			webSocket: &webSocketService{
				client: supervisorClient,
			},
		},
	}
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
