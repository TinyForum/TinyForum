package vo

import (
	"time"
)

type AuthResultVO struct {
	Token          string          `json:"token"`
	User           *UserPrivateVO  `json:"user"`
	DeletionStatus *DeletionStatus `json:"deletion_status,omitempty"`
}

type DeletionStatus struct {
	IsDeleted     bool       `json:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CanRestore    bool       `json:"can_restore"`
	RemainingDays int        `json:"remaining_days,omitempty"`
}
