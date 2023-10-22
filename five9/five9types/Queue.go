package five9types

type QueueID string

type SkillInfo struct {
	ID   QueueID `json:"id"`
	Name string  `json:"name"`
}
