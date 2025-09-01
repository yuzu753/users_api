package domain

import "time"

type UserID int

type User struct {
	ID                 UserID     `json:"id"`
	UserName           string     `json:"user_name"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"-"`
	OtpSecretKey       *string    `json:"-"`
	Type               int        `json:"type"`
	AuthorityData      *string    `json:"authority_data"`
	FailedCount        int        `json:"failed_count"`
	UnlockAt           *time.Time `json:"unlock_at"`
	IsResetPassword    bool       `json:"is_reset_password"`
	LastUpdatePassword time.Time  `json:"last_update_password"`
}