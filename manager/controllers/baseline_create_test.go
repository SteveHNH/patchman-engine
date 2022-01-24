package controllers

import (
	"app/base/core"
	"app/base/database"
	"app/base/utils"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateBaseline(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{
		"name": "my_baseline",
		"inventory_ids": [
			"00000000-0000-0000-0000-000000000005",
			"00000000-0000-0000-0000-000000000006"
		],
        "config": {"to_time": "2022-12-31T12:00:00-04:00"}
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var baselineID int
	ParseResponseBody(t, w.Body.Bytes(), &baselineID)
	database.CheckBaseline(t, baselineID, []string{
		"00000000-0000-0000-0000-000000000005",
		"00000000-0000-0000-0000-000000000006",
	}, `{"to_time": "2022-12-31T12:00:00-04:00"}`, "my_baseline")
	database.DeleteBaseline(t, baselineID)
}

func TestCreateBaselineNameOnly(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{"name": "my_empty_baseline"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var baselineID int
	ParseResponseBody(t, w.Body.Bytes(), &baselineID)
	database.CheckBaseline(t, baselineID, []string{}, "", "my_empty_baseline")
	database.DeleteBaseline(t, baselineID)
}

func TestCreateBaselineNameEmptyString(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{"name": ""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var baselineID int
	ParseResponseBody(t, w.Body.Bytes(), &baselineID)
	database.CheckBaseline(t, baselineID, []string{}, "", "")
	database.DeleteBaseline(t, baselineID)
}

func TestCreateBaselineMissingName(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, "missing required parameter 'name'", errResp.Error)
}

func TestCreateBaselineInvalidRequest(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()
	data := `{"name": 0}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", bytes.NewBufferString(data))
	core.InitRouterWithParams(CreateBaselineHandler, 1, "PUT", "/").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseResponseBody(t, w.Body.Bytes(), &errResp)
	assert.True(t, strings.Contains(errResp.Error,
		"cannot unmarshal number into Go struct field CreateBaselineRequest.name of type string"))
}