package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kramllih/filterService/config"
	"github.com/kramllih/filterService/internal/database"
	_ "github.com/kramllih/filterService/internal/database/mockdb"
	"github.com/kramllih/filterService/internal/httpClient"
	"github.com/kramllih/filterService/internal/logger"
	"github.com/stretchr/testify/assert"
)

func mockController() *Controller {

	return &Controller{
		log:        logger.NewLogger("controller"),
		httpClient: httpClient.MockHTTP(),
	}
}

type MockTransport struct {
	Response    *http.Response
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST" // or PUT
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func TestValidateSimple(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	message := database.Message{
		ID: "1",
		Body: `# Simple Message


This is a simple message.`,
	}

	MockJsonPost(ctx, message)

	cfg := config.RawConfig{
		"database": map[string]interface{}{
			"mockDB": map[string]interface{}{
				"test": "test",
			},
		},
	}

	databaseCfg, err := config.UnpackNamespace("database", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.Load(&databaseCfg)
	if err != nil {
		t.Fatal(err)
	}

	mocktrans := MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"updated":"2022-06-15T19:17:58.3303721Z","words":["harmony","revolutionary","bounce","clue","auction","crew","question","flower","rescue","affair","think","night","morale","route","regular","veil","ensure","communication","undertake","gear","professional","judgment","adult","jaw","death","sex"]}`)),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }

	ctrl := mockController()

	ctrl.httpClient.SetTransport(&mocktrans)

	ctrl.DB = db

	ctrl.Validate(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)

	result := map[string]interface{}{}

	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "your message has been stored.", result["status"])

}

func TestValidateInternalLink(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	message := database.Message{
		ID: "1",
		Body: `# Link Message

- [Section 1](#section1)

## Section 1

This is a simple message.`,
	}

	MockJsonPost(ctx, message)

	cfg := config.RawConfig{
		"database": map[string]interface{}{
			"mockDB": map[string]interface{}{
				"test": "test",
			},
		},
	}

	databaseCfg, err := config.UnpackNamespace("database", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.Load(&databaseCfg)
	if err != nil {
		t.Fatal(err)
	}

	mocktrans := MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"updated":"2022-06-15T19:17:58.3303721Z","words":["harmony","revolutionary","bounce","clue","auction","crew","question","flower","rescue","affair","think","night","morale","route","regular","veil","ensure","communication","undertake","gear","professional","judgment","adult","jaw","death","sex"]}`)),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }

	ctrl := mockController()

	ctrl.httpClient.SetTransport(&mocktrans)

	ctrl.DB = db

	ctrl.Validate(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)

	result := map[string]interface{}{}

	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "your message has been stored.", result["status"])

}

func TestValidateImageLink(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	message := database.Message{
		ID: "1",
		Body: `# Image Message

![The Eiffel Tower](https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/240px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg)`,
	}

	MockJsonPost(ctx, message)

	cfg := config.RawConfig{
		"database": map[string]interface{}{
			"mockDB": map[string]interface{}{
				"test": "test",
			},
		},
	}

	databaseCfg, err := config.UnpackNamespace("database", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.Load(&databaseCfg)
	if err != nil {
		t.Fatal(err)
	}

	mocktrans := MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"updated":"2022-06-15T19:17:58.3303721Z","words":["harmony","revolutionary","bounce","clue","auction","crew","question","flower","rescue","affair","think","night","morale","route","regular","veil","ensure","communication","undertake","gear","professional","judgment","adult","jaw","death","sex"]}`)),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }

	ctrl := mockController()

	ctrl.httpClient.SetTransport(&mocktrans)

	ctrl.DB = db

	ctrl.Validate(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)

	result := map[string]interface{}{}

	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "your message is awaiting approval as it contains image links.", result["status"])

}

func TestValidateExternalLink(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	message := database.Message{
		ID: "1",
		Body: `# Image Message

![google](https://www.google.com)`,
	}

	MockJsonPost(ctx, message)

	cfg := config.RawConfig{
		"database": map[string]interface{}{
			"mockDB": map[string]interface{}{
				"test": "test",
			},
		},
	}

	databaseCfg, err := config.UnpackNamespace("database", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.Load(&databaseCfg)
	if err != nil {
		t.Fatal(err)
	}

	mocktrans := MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"updated":"2022-06-15T19:17:58.3303721Z","words":["harmony","revolutionary","bounce","clue","auction","crew","question","flower","rescue","affair","think","night","morale","route","regular","veil","ensure","communication","undertake","gear","professional","judgment","adult","jaw","death","sex"]}`)),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }

	ctrl := mockController()

	ctrl.httpClient.SetTransport(&mocktrans)

	ctrl.DB = db

	ctrl.Validate(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)

	result := map[string]interface{}{}

	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "your message has has been rejected.", result["status"])

}

func TestValidateBanned(t *testing.T) {

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	message := database.Message{
		ID: "1",
		Body: `# Rejected Language


This message contains adult content`,
	}

	MockJsonPost(ctx, message)

	cfg := config.RawConfig{
		"database": map[string]interface{}{
			"mockDB": map[string]interface{}{
				"test": "test",
			},
		},
	}

	databaseCfg, err := config.UnpackNamespace("database", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.Load(&databaseCfg)
	if err != nil {
		t.Fatal(err)
	}

	mocktrans := MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"updated":"2022-06-15T19:17:58.3303721Z","words":["harmony","revolutionary","bounce","clue","auction","crew","question","flower","rescue","affair","think","night","morale","route","regular","veil","ensure","communication","undertake","gear","professional","judgment","adult","jaw","death","sex"]}`)),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }

	ctrl := mockController()

	ctrl.httpClient.SetTransport(&mocktrans)

	ctrl.DB = db

	ctrl.Validate(ctx)
	assert.EqualValues(t, http.StatusOK, w.Code)

	result := map[string]interface{}{}

	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "your message has has been rejected.", result["status"])

}
