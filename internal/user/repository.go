package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user User) (*User, error)
}

type repository struct {
	db  *sql.DB
	psq sq.StatementBuilderType
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db:  db,
		psq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

func (repo *repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := repo.psq.Select("id", "username", "password", "role", "created_at").From("users").Where(sq.Eq{"username": username}).Limit(1)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRowContext(ctx, sqlStr, args...)

	var user User
	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *repository) Create(ctx context.Context, user User) (*User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := repo.psq.Insert("users").Columns("username", "password", "role", "created_at", "updated_at").Values(user.Username, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	result, err := repo.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id
	return &user, nil
}
