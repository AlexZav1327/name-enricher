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
	updateUserQuery = `
	UPDATE enriched_user
	SET age = $2, gender = $3, country = $4
	WHERE name = $1
	RETURNING name, age, gender, country;
	`
	deleteUserQuery = `
	DELETE FROM enriched_user
	WHERE name = $1;
	`
	name    = "name"
	age     = "age"
	gender  = "gender"
	country = "country"
)

var ErrUserNotFound = errors.New("no such user")

func (p *Postgres) SaveUser(ctx context.Context, user models.ResponseEnrich) error {
	_, err := p.db.Exec(ctx, saveUserQuery, user.Name, user.Age, user.Gender, user.Country)
	if err != nil {
		return fmt.Errorf("p.db.Exec(ctx, saveUserQuery, user.Name, user.Age, user.Gender, user.Country): %w", err)
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

		return models.ResponseEnrich{}, fmt.Errorf("row.Scan(&user.Name, &user.Age, &user.Gender, &user.Country): %w", err)
	}

	return user, nil
}

func (p *Postgres) GetUsersList(ctx context.Context, params models.ListingQueryParams) (
	[]models.ResponseEnrich, error,
) {
	tableColumnsList := map[string]string{
		name:    name,
		age:     age,
		gender:  gender,
		country: country,
	}

	var args []interface{}

	query := `
	SELECT name, age, gender, country
	FROM enriched_user
	WHERE TRUE
	`

	updatedQuery, updatedArgs := p.buildQueryAndArgs(tableColumnsList, args, query, params)

	rows, err := p.db.Query(ctx, updatedQuery, updatedArgs...)
	if err != nil {
		return nil, fmt.Errorf("p.db.Query(ctx, updatedQuery, updatedArgs...): %w", err)
	}

	defer rows.Close()

	usersList := make([]models.ResponseEnrich, 0)

	for rows.Next() {
		var user models.ResponseEnrich

		err = rows.Scan(&user.Name, &user.Age, &user.Gender, &user.Country)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan(&user.Name, &user.Age, &user.Gender, &user.Country): %w", err)
		}

		usersList = append(usersList, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return usersList, nil
}

func (p *Postgres) UpdateUser(ctx context.Context, user models.ResponseEnrich) (models.ResponseEnrich, error) {
	row := p.db.QueryRow(ctx, updateUserQuery, user.Name, user.Age, user.Gender, user.Country)

	var updatedUser models.ResponseEnrich

	err := row.Scan(&updatedUser.Name, &updatedUser.Age, &updatedUser.Gender, &updatedUser.Country)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ResponseEnrich{}, ErrUserNotFound
		}

		return models.ResponseEnrich{}, fmt.Errorf(
			"row.Scan(&updatedUser.Name, &updatedUser.Age, &updatedUser.Gender, &updatedUser.Country): %w", err)
	}

	return updatedUser, nil
}

func (p *Postgres) DeleteUser(ctx context.Context, userName string) error {
	commandTag, err := p.db.Exec(ctx, deleteUserQuery, userName)
	if err != nil {
		return fmt.Errorf("p.db.Exec(ctx, deleteUserQuery, userName): %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return ErrUserNotFound
	}

	return nil
}

func (*Postgres) buildQueryAndArgs(tableColumnsList map[string]string, args []interface{}, query string,
	params models.ListingQueryParams,
) (string, []interface{}) {
	if params.TextFilter != "" {
		args = append(args, "%"+params.TextFilter+"%")
		query += fmt.Sprintf(` AND (name ILIKE $%d OR gender ILIKE $%d OR country ILIKE $%d)`, len(args),
			len(args), len(args))
	}

	order := ` ORDER BY name`

	sorting, ok := tableColumnsList[params.Sorting]
	if ok {
		order = fmt.Sprintf(` ORDER BY %s`, sorting)
	}

	if params.Descending {
		order += ` DESC`
	}

	query += order

	args = append(args, params.ItemsPerPage)
	query += fmt.Sprintf(` LIMIT $%d`, len(args))
	args = append(args, params.Offset)
	query += fmt.Sprintf(` OFFSET $%d`, len(args))

	return query, args
}
