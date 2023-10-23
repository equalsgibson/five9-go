package five9types

type UserFullStateInfo struct {
	CurrentState          Presence `json:"currentState"`
	CurrentStateTime      uint64   `json:"currentStateLong"`
	IsGracefulModeOn      bool     `json:"isGracefulModeOn"`
	PendingState          Presence `json:"pendingState"`
	PendingStateDelayTime uint64   `json:"pendingStateDelayTime"`
}
