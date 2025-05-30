package stores

import (
	"context"

	"github.com/lulzshadowwalker/green-backend/internal"
	"github.com/lulzshadowwalker/green-backend/internal/psql/db"
)

type Users struct {
	q *db.Queries
}

func NewUsers(q *db.Queries) *Users {
	return &Users{q: q}
}

func (u *Users) GetUserByUsername(ctx context.Context, username string) (internal.User, error) {
	user, err := u.q.GetUserByUsername(ctx, username)
	if err != nil {
		return internal.User{}, err
	}
	return internal.User{
		ID:           int(user.ID),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, nil
}
