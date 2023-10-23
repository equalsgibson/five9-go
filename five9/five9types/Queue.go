package five9types

type QueueID string

type QueueInfo struct {
	ID   QueueID `json:"id"`
	Name string  `json:"name"`
}
