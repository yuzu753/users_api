package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, tenantID string, u *User) error
	FindByID(ctx context.Context, tenantID string, id UserID) (*User, error)
	Update(ctx context.Context, tenantID string, u *User) error
	Delete(ctx context.Context, tenantID string, id UserID) error
	Search(ctx context.Context, tenantID string, userName, email string, userType *int, limit, offset int) ([]*User, int, error)
}