package five9

import (
	"encoding/json"
	"errors"
	"fmt"
)

type five9Error struct {
	Five9ExceptionDetail five9ExceptionDetail `json:"five9ExceptionDetail"`
}

type five9ExceptionDetail struct {
	Timestamp uint         `json:"timestamp"`
	ErrorCode int          `json:"errorCode"`
	Message   string       `json:"message"`
	Context   errorContext `json:"context"`
}

type errorContext struct {
	ContextCode string `json:"contextCode"`
	ObjectID    string `json:"objectId"`
}

type Error struct {
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body"`
	Message    string `json:"message"`
}

func (err *Error) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("Five9 REST API Error, Status Code: %d", err.StatusCode)
	}

	return err.Message
}

func (err *Error) UnmarshalJSON(b []byte) error {
	target := five9Error{}
	if targetErr := json.Unmarshal(b, &target); targetErr == nil {
		if target.Five9ExceptionDetail.Message != "" {
			err.Message = target.Five9ExceptionDetail.Message

			return nil
		}
	}

	return nil
}

var (
	ErrWebSocketCacheNotReady error = errors.New("webSocket cache is not ready")
	ErrUnknownUserID          error = errors.New("unknown userID provided")
)
