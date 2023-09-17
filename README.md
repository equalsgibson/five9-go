# five9-go

Five9 Client Library for Go

## Documentation Link

https://webapps.five9.com/assets/files/for_customers/documentation/apis/vcc-agent+supervisor-rest-api-reference-guide.pdf

## Basic Sequence Diagram for connecting to Websocket connection

```mermaid
sequenceDiagram
    participant five9-REST-API
    participant five9-go
    participant five9-Websocket

    five9-go->>+five9-REST-API: POST {REST_API_URL}/appsvcs/rs/svc/auth/login
    note over five9-go, five9-REST-API: Authenticate with the REST API to obtain Authentication Cookies
    five9-REST-API->>-five9-go: Authentication Response
    note over five9-go: Set Authentication Cookies from response
    note over five9-go: Set Active_Datacenter_URL from response
    five9-go->>+five9-REST-API: GET {Active_Datacenter_URL}/supsvcs/rs/svc/auth/metadata
    note over five9-go, five9-REST-API: Send request to Metadata endpoint to validate session
    five9-REST-API->>-five9-go: Metadata Response
    five9-go->five9-Websocket: WSS Connection to wss://{Active_Datacenter_URL}/supsvcs/sws/ws
    five9-Websocket->>five9-go: Receive [EventID 1010] (Server Connected) Message
    five9-Websocket-->>five9-go: Receive Supervisor Websocket Messages
    loop Every 15 seconds
    note over five9-go, five9-Websocket: Keep WSS Connection alive with frequent ping messages.
        five9-go-->>five9-Websocket: Ping WSS Connection
        five9-go->>five9-go: Increment ping count by 1
        alt Ping count is greater than 2
            five9-go->>five9-Websocket: Assume connection is in error state. Close WSS connection.
        end
        five9-Websocket-->>five9-go: Receive Pong [EventID 1202 (Pong)] Message
         five9-go->>five9-go: Reset ping count to 0
    end;
```
