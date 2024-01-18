//nolint:wrapcheck
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	pg              store
	ageResolver     AgeResolver
	genderResolver  GenderResolver
	countryResolver CountryResolver
	log             *logrus.Entry
	metrics         *metrics
}

func New(pg store, age AgeResolver, gender GenderResolver, country CountryResolver, log *logrus.Logger) *Service {
	return &Service{
		pg:              pg,
		ageResolver:     age,
		genderResolver:  gender,
		countryResolver: country,
		log:             log.WithField("module", "service"),
		metrics:         newMetrics(),
	}
}

type store interface {
	GetUser(ctx context.Context, userName string) (models.ResponseEnrich, error)
	GetUsersList(ctx context.Context, params models.ListingQueryParams) ([]models.ResponseEnrich, error)
	SaveUser(ctx context.Context, user models.ResponseEnrich) error
	UpdateUser(ctx context.Context, user models.ResponseEnrich) (models.ResponseEnrich, error)
	DeleteUser(ctx context.Context, userName string) error
}

type AgeResolver interface {
	GetAge(ctx context.Context, name string) (int, error)
}

type GenderResolver interface {
	GetGender(ctx context.Context, name string) (string, error)
}

type CountryResolver interface {
	GetCountry(ctx context.Context, name string) (string, error)
}

func (s *Service) EnrichUser(ctx context.Context, userName models.RequestEnrich) (models.ResponseEnrich, error) {
	userNameEnriched, err := s.pg.GetUser(ctx, userName.Name)
	if err == nil && userNameEnriched.Name == userName.Name && userNameEnriched.Surname == userName.Surname &&
		userNameEnriched.Patronymic == userName.Patronymic {
		return userNameEnriched, nil
	}

	userNameEnriched = models.ResponseEnrich{
		RequestEnrich: userName,
	}

	eg, egCtx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		age, err := s.ageResolver.GetAge(egCtx, userName.Name)
		if err != nil {
			return err
		}
		userNameEnriched.Age = age

		return nil
	})

	eg.Go(func() error {
		gender, err := s.genderResolver.GetGender(egCtx, userName.Name)
		if err != nil {
			return err
		}
		userNameEnriched.Gender = gender

		return nil
	})

	eg.Go(func() error {
		country, err := s.countryResolver.GetCountry(egCtx, userName.Name)
		if err != nil {
			return err
		}
		userNameEnriched.Country = country

		return nil
	})

	err = eg.Wait()
	if err != nil {
		return models.ResponseEnrich{}, fmt.Errorf("eg.Wait(): %w", err)
	}

	started := time.Now()
	defer func() {
		s.metrics.duration.WithLabelValues("save_user").Observe(time.Since(started).Seconds())
	}()

	err = s.pg.SaveUser(ctx, userNameEnriched)
	if err != nil {
		return userNameEnriched, fmt.Errorf("s.pg.SaveUser(ctx, userNameEnriched): %w", err)
	}

	return userNameEnriched, nil
}

func (s *Service) GetUsersList(ctx context.Context, params models.ListingQueryParams) ([]models.ResponseEnrich, error) {
	usersList, err := s.pg.GetUsersList(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.pg.GetUsersList(ctx, params): %w", err)
	}

	return usersList, nil
}

func (s *Service) UpdateUser(ctx context.Context, user models.ResponseEnrich) (models.ResponseEnrich, error) {
	currentUser, err := s.pg.GetUser(ctx, user.Name)
	if err != nil {
		return models.ResponseEnrich{}, fmt.Errorf("s.pg.GetUser(ctx, user.Name): %w", err)
	}

	if user.Surname == "" {
		user.Surname = currentUser.Surname
	}

	if user.Patronymic == "" {
		user.Patronymic = currentUser.Patronymic
	}

	if user.Age == 0 {
		user.Age = currentUser.Age
	}

	if user.Gender == "" {
		user.Gender = currentUser.Gender
	}

	if user.Country == "" {
		user.Country = currentUser.Country
	}

	started := time.Now()
	defer func() {
		s.metrics.duration.WithLabelValues("update_user").Observe(time.Since(started).Seconds())
	}()

	updatedUser, err := s.pg.UpdateUser(ctx, user)
	if err != nil {
		return models.ResponseEnrich{}, fmt.Errorf("s.pg.UpdateUser(ctx, user): %w", err)
	}

	return updatedUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, userName string) error {
	started := time.Now()
	defer func() {
		s.metrics.duration.WithLabelValues("delete_user").Observe(time.Since(started).Seconds())
	}()

	err := s.pg.DeleteUser(ctx, userName)
	if err != nil {
		return fmt.Errorf("s.pg.DeleteUser(ctx, userName): %w", err)
	}

	return nil
}
