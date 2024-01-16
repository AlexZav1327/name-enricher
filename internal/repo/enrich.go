package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
)

type Enrich struct {
	log *logrus.Entry
}

func New(log *logrus.Logger) *Enrich {
	return &Enrich{
		log: log.WithField("module", "enrich"),
	}
}

func (e *Enrich) EnrichAge(ctx context.Context, name string) (int, error) {
	endpoint := fmt.Sprintf("https://api.agify.io/?name=%s", name)

	var respData models.AgeEnriched

	err := e.sendRequest(ctx, endpoint, &respData)
	if err != nil {
		return 0, fmt.Errorf("sendRequest: %w", err)
	}

	return respData.Age, nil
}

func (e *Enrich) EnrichGender(ctx context.Context, name string) (string, error) {
	endpoint := fmt.Sprintf("https://api.genderize.io/?name=%s", name)

	var respData models.GenderEnricher

	err := e.sendRequest(ctx, endpoint, &respData)
	if err != nil {
		return "", fmt.Errorf("sendRequest: %w", err)
	}

	return respData.Gender, nil
}

func (e *Enrich) EnrichCountry(ctx context.Context, name string) (string, error) {
	endpoint := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)

	var respData models.CountryEnrichedList

	err := e.sendRequest(ctx, endpoint, &respData)
	if err != nil {
		return "", fmt.Errorf("sendRequest: %w", err)
	}

	return respData.Country[0].CountryID, nil
}

func (e *Enrich) sendRequest(ctx context.Context, endpoint string, respData interface{}) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do: %w", err)
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			e.log.Warningf("response.Body.Close: %s", err)
		}
	}()

	err = json.NewDecoder(response.Body).Decode(&respData)
	if err != nil {
		return fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	return nil
}
