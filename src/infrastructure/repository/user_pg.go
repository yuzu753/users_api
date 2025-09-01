package repository

import (
	"context"
	"fmt"
	"strings"
	"users_api/src/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPg struct {
	db *pgxpool.Pool
}

func NewUserPg(db *pgxpool.Pool) domain.UserRepository {
	return &UserPg{db: db}
}

func (r *UserPg) Create(ctx context.Context, tenantID string, u *domain.User) error {
	table := fmt.Sprintf("users_%s", tenantID)
	query := fmt.Sprintf(`
		INSERT INTO %s (user_name, email, password_hash, otp_secret_key, type, authority_data, failed_count, unlock_at, is_reset_password, last_update_password)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, table)
	
	return r.db.QueryRow(ctx, query, u.UserName, u.Email, u.PasswordHash, u.OtpSecretKey, u.Type, u.AuthorityData, u.FailedCount, u.UnlockAt, u.IsResetPassword, u.LastUpdatePassword).Scan(&u.ID)
}

func (r *UserPg) FindByID(ctx context.Context, tenantID string, id domain.UserID) (*domain.User, error) {
	table := fmt.Sprintf("users_%s", tenantID)
	query := fmt.Sprintf(`
		SELECT id, user_name, email, password_hash, otp_secret_key, type, authority_data, failed_count, unlock_at, is_reset_password, last_update_password
		FROM %s WHERE id = $1
	`, table)
	
	var u domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.UserName, &u.Email, &u.PasswordHash, &u.OtpSecretKey, &u.Type, &u.AuthorityData, &u.FailedCount, &u.UnlockAt, &u.IsResetPassword, &u.LastUpdatePassword,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserPg) Update(ctx context.Context, tenantID string, u *domain.User) error {
	table := fmt.Sprintf("users_%s", tenantID)
	query := fmt.Sprintf(`
		UPDATE %s SET user_name = $2, email = $3, password_hash = $4, otp_secret_key = $5, type = $6, authority_data = $7, failed_count = $8, unlock_at = $9, is_reset_password = $10, last_update_password = $11
		WHERE id = $1
	`, table)
	
	_, err := r.db.Exec(ctx, query, u.ID, u.UserName, u.Email, u.PasswordHash, u.OtpSecretKey, u.Type, u.AuthorityData, u.FailedCount, u.UnlockAt, u.IsResetPassword, u.LastUpdatePassword)
	return err
}

func (r *UserPg) Delete(ctx context.Context, tenantID string, id domain.UserID) error {
	table := fmt.Sprintf("users_%s", tenantID)
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *UserPg) Search(ctx context.Context, tenantID, userName, email string, userType *int, limit, offset int) ([]*domain.User, int, error) {
	table := fmt.Sprintf("users_%s", tenantID)
	where := []string{"1=1"}
	args := []any{}

	if userName != "" {
		where = append(where, fmt.Sprintf("user_name ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, userName)
	}
	if email != "" {
		where = append(where, fmt.Sprintf("email ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, email)
	}
	if userType != nil {
		where = append(where, fmt.Sprintf("type = $%d", len(args)+1))
		args = append(args, *userType)
	}

	countSQL := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s", table, strings.Join(where, " AND "))
	var total int
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf(`
		SELECT id, user_name, email, type, authority_data, failed_count, unlock_at, is_reset_password, last_update_password 
		FROM %s WHERE %s LIMIT $%d OFFSET $%d
	`, table, strings.Join(where, " AND "), len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.UserName, &u.Email, &u.Type, &u.AuthorityData, &u.FailedCount, &u.UnlockAt, &u.IsResetPassword, &u.LastUpdatePassword); err != nil {
			return nil, 0, err
		}
		items = append(items, &u)
	}
	
	return items, total, nil
}