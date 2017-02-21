package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *gin.Engine = getMainEngine()

func TestPingPong(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ping", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("Response code should be Ok, was: %d", res.Code)
	}

	bodyAsString := res.Body.String()

	if bodyAsString != "pong" {
		t.Errorf("Response body should be `pongd`, was  %s", bodyAsString)
	}
}

func TestHelloWorld(t *testing.T) {
	req, _ := http.NewRequest("GET", "/hello", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("Response code should be Ok, was: %d", res.Code)
	}

	bodyAsString := res.Body.String()

	if bodyAsString != "Hello, World!" {
		t.Errorf("Response body should be `Hello, World`, was  %s", bodyAsString)
	}
}
