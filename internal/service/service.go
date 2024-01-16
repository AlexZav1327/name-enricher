package service

import (
	"context"
	"fmt"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
)

type Service struct {
	er  Enricher
	log *logrus.Entry
}

type Enricher interface {
	EnrichAge(ctx context.Context, name string) (int, error)
	EnrichGender(ctx context.Context, name string) (string, error)
	EnrichCountry(ctx context.Context, name string) (string, error)
}

func New(er Enricher, log *logrus.Logger) *Service {
	return &Service{
		er:  er,
		log: log.WithField("module", "service"),
	}
}

func (s *Service) Enrich(ctx context.Context, personName models.RequestEnrich) (models.ResponseEnrich, error) {
	var responseEnrich models.ResponseEnrich

	responseEnrich.Name = personName.Name
	responseEnrich.Surname = personName.Surname
	responseEnrich.Patronymic = personName.Patronymic

	age, err := s.GetAge(ctx, personName.Name)
	if err != nil {
		return models.ResponseEnrich{}, fmt.Errorf("GetAge: %w", err)
	}

	gender, err := s.GetGender(ctx, personName.Name)
	if err != nil {
		return models.ResponseEnrich{}, fmt.Errorf("GetGender: %w", err)
	}

	country, err := s.GetCountry(ctx, personName.Name)
	if err != nil {
		return models.ResponseEnrich{}, fmt.Errorf("GetCountry: %w", err)
	}

	responseEnrich.Age = age
	responseEnrich.Gender = gender
	responseEnrich.Country = country

	return responseEnrich, nil
}

func (s *Service) GetAge(ctx context.Context, name string) (int, error) {
	age, err := s.er.EnrichAge(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("EnrichAge: %w", err)
	}

	return age, nil
}

func (s *Service) GetGender(ctx context.Context, name string) (string, error) {
	gender, err := s.er.EnrichGender(ctx, name)
	if err != nil {
		return "", fmt.Errorf("EnrichGender: %w", err)
	}

	return gender, nil
}

func (s *Service) GetCountry(ctx context.Context, name string) (string, error) {
	country, err := s.er.EnrichCountry(ctx, name)
	if err != nil {
		return "", fmt.Errorf("EnrichCountry: %w", err)
	}

	return country, nil
}
