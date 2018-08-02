package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/fraenky8/zlr-ca/pkg/domain"
	"github.com/fraenky8/zlr-ca/pkg/mock"
	"github.com/fraenky8/zlr-ca/pkg/storage/repos"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	basicAuthHeaderFrank = `Basic ZnJhbms6ZnI0bmsh`

	requestContentType  = "application/json"
	responseContentType = "application/json; charset=utf-8"

	icecreamProductId1 = "602"
	icecreamProductId2 = "610"

	icecream = `
		[{
            "productId": "602",
            "name": "Banana Split",
            "description": "Banana & Strawberry Ice Creams with Walnuts, Fudge Chunks & a Fudge Swirl",
            "story": "We turned the classic ice cream parlor sundae you've always loved into the at-home flavor creation you've always wanted. Enjoy!",
            "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing.png",
            "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing-open.png",
            "allergy_info": "may contain other tree nuts",
            "dietary_certifications": "Kosher"
        }]
	`

	missingProductIdIcecream = `
		[{
            "name": "Banana Split",
            "description": "Banana & Strawberry Ice Creams with Walnuts, Fudge Chunks & a Fudge Swirl",
            "story": "We turned the classic ice cream parlor sundae you've always loved into the at-home flavor creation you've always wanted. Enjoy!",
            "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing.png",
            "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing-open.png",
            "allergy_info": "may contain other tree nuts",
            "dietary_certifications": "Kosher"
        }]
	`

	missingNameIcecream = `
		[{
			"productId": "602",
            "description": "Banana & Strawberry Ice Creams with Walnuts, Fudge Chunks & a Fudge Swirl",
            "story": "We turned the classic ice cream parlor sundae you've always loved into the at-home flavor creation you've always wanted. Enjoy!",
            "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing.png",
            "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing-open.png",
            "allergy_info": "may contain other tree nuts",
            "dietary_certifications": "Kosher"
        }]
	`

	faultyProductIdIcecream = `
		[{
            "productId": "60xx2",
            "name": "Banana Split",
            "description": "Banana & Strawberry Ice Creams with Walnuts, Fudge Chunks & a Fudge Swirl",
            "story": "We turned the classic ice cream parlor sundae you've always loved into the at-home flavor creation you've always wanted. Enjoy!",
            "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing.png",
            "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/banana-split-landing-open.png",
            "allergy_info": "may contain other tree nuts",
            "dietary_certifications": "Kosher"
        }]
	`
)

func TestReadIcecream_withoutAuthorization_returnsStatusUnauthorized(t *testing.T) {
	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/icecreams/"+icecreamProductId1, nil)
	assert.Nil(t, err)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestReadIcecream_withInvalidAuthorization_returnsStatusUnauthorized(t *testing.T) {
	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/icecreams/"+icecreamProductId1, nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Basic s0m3th!ng")

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestReadSingleIcecream_withValidAuthorization_returnsStatusOkAndSingleIcecreamData(t *testing.T) {
	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	is.ReadsFn = func(ids []int64) ([]*domain.Icecream, error) {
		return []*domain.Icecream{{
			ProductID: icecreamProductId1,
			// ... more data
		}}, nil
	}

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/icecreams/"+icecreamProductId1, nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	response := struct {
		Status string
		Data   struct {
			Icecream struct {
				domain.Icecream
			}
		}
	}{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusOk, response.Status)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, icecreamProductId1, response.Data.Icecream.ProductID)

	assert.True(t, is.ReadsInvoked)
}

func TestReadMultipleIcecream_withValidAuthorization_returnsStatusOkAndMultipleIcecreamData(t *testing.T) {
	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	is.ReadsFn = func(ids []int64) ([]*domain.Icecream, error) {
		return []*domain.Icecream{{
			ProductID: icecreamProductId1,
			// ... more data
		}, {
			ProductID: icecreamProductId2,
			// ... more data
		}}, nil
	}

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/icecreams/"+icecreamProductId1+","+icecreamProductId2, nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	response := struct {
		Status string
		Data   struct {
			Icecreams []struct {
				domain.Icecream
			}
		}
	}{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusOk, response.Status)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 2, len(response.Data.Icecreams))
	assert.Equal(t, icecreamProductId1, response.Data.Icecreams[0].ProductID)
	assert.Equal(t, icecreamProductId2, response.Data.Icecreams[1].ProductID)

	assert.True(t, is.ReadsInvoked)
}

func TestCreateIcecream_withWrongContentType_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)
}

func TestCreateIcecream_withNilBody_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", nil)
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)
}

func TestCreateIcecream_withEmptyBody_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(""))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)
}

func TestCreateIcecream_withFaultyData_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(`{"data":["foo", "bar"]"}`))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)
}

func TestCreateIcecream_withMissingProductId_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(missingProductIdIcecream))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)
}

func TestCreateIcecream_withMissingName_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(missingNameIcecream))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)
}

func TestCreateIcecream_withFaultyProductId_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(faultyProductIdIcecream))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)
}

func TestCreateIcecream_withExistingIcecream_returnsFailResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	is.ReadsFn = func(ids []int64) ([]*domain.Icecream, error) {
		// existing icecream
		var icecreams []*domain.Icecream
		err = json.Unmarshal([]byte(icecream), &icecreams)
		assert.Nil(t, err)

		return icecreams, nil
	}

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(icecream))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusFail, response.Status)
	assert.NotEmpty(t, response.Data)

	assert.True(t, is.ReadsInvoked)
}

func TestCreateIcecream_withNewIcecreamButDatabaseError_returnsErrorResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	is.ReadsFn = func(ids []int64) ([]*domain.Icecream, error) {
		return nil, nil
	}

	is.CreatesFn = func(icecreams []*domain.Icecream) ([]int64, error) {
		// simulation database error
		return nil, fmt.Errorf("foreign key constraint violated")
	}

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(icecream))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusError, response.Status)
	assert.Empty(t, response.Data)
	assert.NotEmpty(t, response.Message)

	assert.True(t, is.ReadsInvoked)
	assert.True(t, is.CreatesInvoked)
}

func TestCreateIcecream_withNewValidIcecream_returnsSuccessResponse(t *testing.T) {

	// given
	is := &mock.IcecreamService{}

	s, err := NewServer(
		&ServerConfig{Mode: gin.ReleaseMode},
		&repos.Repository{
			IcecreamService:                  is,
			IngredientService:                &mock.IngredientService{},
			SourcingValueService:             &mock.SourcingValueService{},
			IcecreamHasIngredientsService:    &mock.IcecreamHasIngredientsService{},
			IcecreamHasSourcingValuesService: &mock.IcecreamHasSourcingValuesService{},
		},
	)
	assert.Nil(t, err)

	is.ReadsFn = func(ids []int64) ([]*domain.Icecream, error) {
		return nil, nil
	}

	is.CreatesFn = func(icecreams []*domain.Icecream) ([]int64, error) {
		id, _ := strconv.Atoi(icecreamProductId1)
		return []int64{int64(id)}, nil
	}

	// when
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/icecreams", strings.NewReader(icecream))
	r.Header.Set("Content-Type", requestContentType)
	r.Header.Set("Authorization", basicAuthHeaderFrank)

	s.ServeHTTP(w, r)

	// then
	assert.Equal(t, responseContentType, w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, http.StatusCreated, w.Code)

	response := struct {
		Status string
		Data   struct {
			Icecreams []struct {
				domain.Icecream
			}
		}
	}{}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, StatusOk, response.Status)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, icecreamProductId1, response.Data.Icecreams[0].ProductID)

	assert.True(t, is.ReadsInvoked)
	assert.True(t, is.CreatesInvoked)
}
