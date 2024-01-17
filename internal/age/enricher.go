package age

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlexZav1327/name-enricher/internal/models"
	"github.com/sirupsen/logrus"
)

type Age struct {
	log     *logrus.Entry
	metrics *metrics
}

func New(log *logrus.Logger) *Age {
	return &Age{
		log:     log.WithField("module", "age"),
		metrics: newMetrics(),
	}
}

func (a *Age) GetAge(ctx context.Context, name string) (int, error) {
	endpoint := fmt.Sprintf("https://api.agify.io/?name=%s", name)

	var respData models.AgeEnriched

	err := a.sendRequest(ctx, endpoint, &respData)
	if err != nil {
		return 0, fmt.Errorf("a.sendRequest(ctx, endpoint, &respData): %w", err)
	}

	if respData.Age == 0 {
		return 0, models.ErrNameNotValid
	}

	return respData.Age, nil
}

func (a *Age) sendRequest(ctx context.Context, endpoint string, respData interface{}) error {
	started := time.Now()
	defer func() {
		a.metrics.duration.Observe(time.Since(started).Seconds())
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
			a.log.Warningf("response.Body.Close(): %s", err)
		}
	}()

	err = json.NewDecoder(response.Body).Decode(&respData)
	if err != nil {
		return fmt.Errorf("json.NewDecoder(response.Body).Decode(&respData): %w", err)
	}

	return nil
}
