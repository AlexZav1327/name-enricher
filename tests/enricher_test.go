package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AlexZav1327/name-enricher/internal/models"
)

func (s *IntegrationTestSuite) TestServiceCRUD() {
	s.Run("enrich user name with details normal case", func() {
		ctx := context.Background()

		req := models.RequestEnrich{
			Name: "Liza",
		}

		var respData models.ResponseEnrich

		resp := s.sendRequest(ctx, http.MethodPost, url+enrichNameEndpoint, req, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(req.Name, respData.Name)

		var respAge models.AgeEnriched

		endpoint := fmt.Sprintf("https://api.agify.io/?name=%s", req.Name)
		_ = s.sendRequest(ctx, http.MethodGet, endpoint, nil, &respAge)

		s.Require().Equal(respData.Age, respAge.Age)

		var respGender models.GenderEnricher

		endpoint = fmt.Sprintf("https://api.genderize.io/?name=%s", req.Name)
		_ = s.sendRequest(ctx, http.MethodGet, endpoint, nil, &respGender)

		s.Require().Equal(respData.Gender, respGender.Gender)

		var respCountry models.CountryEnrichedList

		endpoint = fmt.Sprintf("https://api.nationalize.io/?name=%s", req.Name)
		_ = s.sendRequest(ctx, http.MethodGet, endpoint, nil, &respCountry)

		s.Require().Equal(respData.Country, respCountry.Country[0].CountryID)
	})
	s.Run("update user normal case", func() {
		ctx := context.Background()

		req := models.RequestEnrich{
			Name: "Lana",
		}

		var respData models.ResponseEnrich

		_ = s.sendRequest(ctx, http.MethodPost, url+enrichNameEndpoint, req, &respData)

		reqUpdate := respData
		reqUpdate.Age = 99
		reqUpdate.Gender = "male"
		reqUpdate.Country = "XY"

		userNameEndpoint := req.Name
		resp := s.sendRequest(ctx, http.MethodPatch, url+updateUserEndpoint+userNameEndpoint, reqUpdate, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(reqUpdate.Age, respData.Age)
		s.Require().Equal(reqUpdate.Gender, respData.Gender)
		s.Require().Equal(reqUpdate.Country, respData.Country)
	})
	s.Run("delete user normal case", func() {
		ctx := context.Background()

		req := models.RequestEnrich{
			Name: "Alex",
		}
		_ = s.sendRequest(ctx, http.MethodPost, url+enrichNameEndpoint, req, nil)

		userNameEndpoint := req.Name
		resp := s.sendRequest(ctx, http.MethodDelete, url+deleteUserEndpoint+userNameEndpoint, nil, nil)

		s.Require().Equal(http.StatusNoContent, resp.StatusCode)
	})
}

func (s *IntegrationTestSuite) TestUsersList() {
	s.Run("get empty list of users normal case", func() {
		ctx := context.Background()

		var respData []models.ResponseEnrich

		resp := s.sendRequest(ctx, http.MethodGet, url+usersListEndpoint, nil, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal([]models.ResponseEnrich{}, respData)
	})
	s.Run("get users list normal case", func() {
		ctx := context.Background()

		req := models.RequestEnrich{
			Name: "Alex",
		}
		_ = s.sendRequest(ctx, http.MethodPost, url+enrichNameEndpoint, req, nil)

		req.Name = "Kate"
		_ = s.sendRequest(ctx, http.MethodPost, url+enrichNameEndpoint, req, nil)

		req.Name = "Liza"
		_ = s.sendRequest(ctx, http.MethodPost, url+enrichNameEndpoint, req, nil)

		var respData []models.ResponseEnrich

		resp := s.sendRequest(ctx, http.MethodGet, url+usersListEndpoint, nil, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(3, len(respData))

		queryParams := "?textFilter=Liza"
		resp = s.sendRequest(ctx, http.MethodGet, url+usersListEndpoint+queryParams, nil, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(1, len(respData))
		s.Require().Equal("Liza", respData[0].Name)

		queryParams = "?sorting=name&descending=true"
		resp = s.sendRequest(ctx, http.MethodGet, url+usersListEndpoint+queryParams, nil, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal("Liza", respData[0].Name)

		queryParams = "?itemsPerPage=2"
		resp = s.sendRequest(ctx, http.MethodGet, url+usersListEndpoint+queryParams, nil, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(2, len(respData))

		queryParams = "?offset=1"
		resp = s.sendRequest(ctx, http.MethodGet, url+usersListEndpoint+queryParams, nil, &respData)

		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(2, len(respData))
	})
}
