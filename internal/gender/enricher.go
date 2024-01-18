package gender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
)

type Gender struct {
	log     *logrus.Entry
	metrics *metrics
}

func New(log *logrus.Logger) *Gender {
	return &Gender{
		log:     log.WithField("module", "gender"),
		metrics: newMetrics(),
	}
}

func (g *Gender) GetGender(ctx context.Context, name string) (string, error) {
	endpoint := fmt.Sprintf("https://api.genderize.io/?name=%s", name)

	var respData models.GenderEnriched

	err := g.sendRequest(ctx, endpoint, &respData)
	if err != nil {
		return "", fmt.Errorf("g.sendRequest(ctx, endpoint, &respData): %w", err)
	}

	if respData.Gender == "" {
		return "", models.ErrNameNotValid
	}

	return respData.Gender, nil
}

func (g *Gender) sendRequest(ctx context.Context, endpoint string, respData interface{}) error {
	started := time.Now()
	defer func() {
		g.metrics.duration.Observe(time.Since(started).Seconds())
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
			g.log.Warningf("response.Body.Close(): %s", err)
		}
	}()

	err = json.NewDecoder(response.Body).Decode(&respData)
	if err != nil {
		return fmt.Errorf("json.NewDecoder(response.Body).Decode(&respData): %w", err)
	}

	return nil
}
