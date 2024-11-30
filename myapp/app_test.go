package myapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPathHandler(t *testing.T) {
	// Code here
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello, World", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	// Code here
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello, World!", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	// Code here
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=seongwon", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello, seongwon!", string(data))
}

func TestFooHandler_WithoutJson(t *testing.T) {
	// Code here
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "Bad Request")
}

func TestFooHandler_WithJson(t *testing.T) {
	// Code here
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo",
		strings.NewReader(`{"first_name":"seongwon","last_name":"jeong","email":"seongwon@naver.com"}`))

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(err)
	assert.Equal("seongwon", user.FirstName)
	assert.Equal("jeong", user.LastName)
	assert.Equal("seongwon@naver.com", user.Email)
}
