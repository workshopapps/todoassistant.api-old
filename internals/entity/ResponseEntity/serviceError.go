package ResponseEntity

import "time"

type ServiceError struct {
	Time        string `json:"time"`
	Description string `json:"description"`
	Error       any    `json:"error,omitempty"`
}

func NewCustomServiceError(description string, error any) *ServiceError {
	return &ServiceError{Time: time.Now().Format(time.RFC3339), Description: description, Error: error}
}

func NewInternalServiceError(error any) *ServiceError {
	return &ServiceError{Time: time.Now().Format(time.RFC3339), Description: "Internal Service Error", Error: error}
}

func NewValidatingError(error any) *ServiceError {
	return &ServiceError{Time: time.Now().Format(time.RFC3339), Description: "BadInput Request", Error: error}
}
