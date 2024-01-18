package country

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
)

type Country struct {
	log     *logrus.Entry
	metrics *metrics
}

func New(log *logrus.Logger) *Country {
	return &Country{
		log:     log.WithField("module", "country"),
		metrics: newMetrics(),
	}
}

func (c *Country) GetCountry(ctx context.Context, name string) (string, error) {
	endpoint := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)

	var respData models.CountryEnrichedList

	err := c.sendRequest(ctx, endpoint, &respData)
	if err != nil {
		return "", fmt.Errorf("c.sendRequest(ctx, endpoint, &respData): %w", err)
	}

	if len(respData.Country) == 0 {
		return "", models.ErrNameNotValid
	}

	return respData.Country[0].CountryID, nil
}

func (c *Country) sendRequest(ctx context.Context, endpoint string, respData interface{}) error {
	started := time.Now()
	defer func() {
		c.metrics.duration.Observe(time.Since(started).Seconds())
	}()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil): %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do(request): %w", err)
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			c.log.Warningf("response.Body.Close(): %s", err)
		}
	}()

	err = json.NewDecoder(response.Body).Decode(&respData)
	if err != nil {
		return fmt.Errorf("json.NewDecoder(response.Body).Decode(&respData): %w", err)
	}

	return nil
}
