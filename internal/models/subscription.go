package models

import "github.com/google/uuid"

type Subscription struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"` // Format: MM-YYYY
	EndDate     *string   `json:"end_date,omitempty"` // Format: MM-YYYY, can be null
}
