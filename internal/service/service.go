//nolint:wrapcheck
package service

import (
	"context"
	"fmt"

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
}

func New(pg store, age AgeResolver, gender GenderResolver, country CountryResolver, log *logrus.Logger) *Service {
	return &Service{
		pg:              pg,
		ageResolver:     age,
		genderResolver:  gender,
		countryResolver: country,
		log:             log.WithField("module", "service"),
	}
}

type store interface {
	GetUser(ctx context.Context, userName string) (models.ResponseEnrich, error)
	SaveUser(ctx context.Context, user models.ResponseEnrich) error
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

func (s *Service) HandleUser(ctx context.Context, userName models.RequestEnrich) (models.ResponseEnrich, error) {
	userNameEnriched, err := s.getUser(ctx, userName)
	if err == nil {
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
		return models.ResponseEnrich{}, fmt.Errorf("eg.Wait: %w", err)
	}

	err = s.pg.SaveUser(ctx, userNameEnriched)
	if err != nil {
		return userNameEnriched, fmt.Errorf("pg.SaveUser: %w", err)
	}

	return userNameEnriched, nil
}

func (s *Service) getUser(ctx context.Context, userName models.RequestEnrich) (models.ResponseEnrich, error) {
	userNameEnriched, err := s.pg.GetUser(ctx, userName.Name)
	if err != nil {
		return models.ResponseEnrich{}, fmt.Errorf("pg.GetUser: %w", err)
	}

	return userNameEnriched, nil
}
