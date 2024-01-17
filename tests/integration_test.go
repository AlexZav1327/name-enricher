package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/AlexZav1327/name-enricher/internal/age"
	"github.com/AlexZav1327/name-enricher/internal/country"
	"github.com/AlexZav1327/name-enricher/internal/gender"
	"github.com/AlexZav1327/name-enricher/internal/storage"

	"github.com/AlexZav1327/name-enricher/internal/server"
	"github.com/AlexZav1327/name-enricher/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const (
	port               = 5005
	host               = ""
	dsn                = "postgres://user:secret@localhost:5436/postgres?sslmode=disable"
	enrichNameEndpoint = "/api/v1/user/enrich"
	updateUserEndpoint = "/api/v1/user/update"
	deleteUserEndpoint = "/api/v1/user/delete"
	usersListEndpoint  = "/api/v1/users"
)

var url = fmt.Sprintf("http://localhost:%d", port)

type IntegrationTestSuite struct {
	suite.Suite
	pg      *storage.Postgres
	server  *server.Server
	service *service.Service
	age     *age.Age
	gender  *gender.Gender
	country *country.Country
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	logger := logrus.StandardLogger()

	var err error

	s.pg, err = storage.ConnectDB(ctx, dsn, logger)
	s.Require().NoError(err)

	err = s.pg.Migrate(migrate.Up)
	s.Require().NoError(err)

	s.age = age.New(logger)
	s.gender = gender.New(logger)
	s.country = country.New(logger)
	s.service = service.New(s.pg, s.age, s.gender, s.country, logger)
	s.server = server.New(host, port, s.service, logger)

	go func() {
		err = s.server.Run(ctx)
		s.Require().NoError(err)
	}()

	time.Sleep(250 * time.Millisecond)
}

func (s *IntegrationTestSuite) TearDownTest() {
	ctx := context.Background()

	err := s.pg.TruncateTable(ctx, "enriched_user")
	s.Require().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) sendRequest(ctx context.Context, method, endpoint string,
	body, dest interface{},
) *http.Response {
	s.T().Helper()

	reqBody, err := json.Marshal(body)
	s.Require().NoError(err)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(reqBody))
	s.Require().NoError(err)

	req.Header.Set("Context-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)

	defer func() {
		err = resp.Body.Close()
		s.Require().NoError(err)
	}()

	if dest != nil {
		err = json.NewDecoder(resp.Body).Decode(&dest)
		s.Require().NoError(err)
	}

	return resp
}
