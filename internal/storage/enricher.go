package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/jackc/pgx/v5"
)

const (
	saveUserQuery = `
	INSERT INTO enriched_user (name, age, gender, country)
	VALUES ($1, $2, $3, $4);
	`
	getUserQuery = `
	SELECT name, age, gender, country
	FROM enriched_user
	WHERE name = $1
	`
)

var ErrUserNotFound = errors.New("no such user")

func (p *Postgres) SaveUser(ctx context.Context, user models.ResponseEnrich) error {
	_, err := p.db.Exec(ctx, saveUserQuery, user.Name, user.Age, user.Gender, user.Country)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (p *Postgres) GetUser(ctx context.Context, userName string) (models.ResponseEnrich, error) {
	row := p.db.QueryRow(ctx, getUserQuery, userName)

	var user models.ResponseEnrich

	err := row.Scan(&user.Name, &user.Age, &user.Gender, &user.Country)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ResponseEnrich{}, ErrUserNotFound
		}

		return models.ResponseEnrich{}, fmt.Errorf("row.Scan: %w", err)
	}

	return user, nil
}
