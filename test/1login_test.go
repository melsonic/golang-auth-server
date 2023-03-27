package test

import (
	"bytes"
	"encoding/json"
	"example/auth/internal/pkg/controller"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	router := gin.Default()

	router.POST("/login", controller.Login)

	requestData := []byte(`{"username": "u3name","password": "u3pword"}`)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestData))

	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP recorder
	respRecorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(respRecorder, req)

	// Verify the response status code
	if status := respRecorder.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	// Parse the response body
	var responseBody map[string]interface{}
	if err := json.Unmarshal(respRecorder.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Retrieve the access token from the response body
	accessToken, ok := responseBody["access_token"].(string)
	if !ok {
		t.Fatalf("Failed to retrieve access token from response body")
	}

	// json response message
	jsonMessage := "user logged in succesfully"

	// Verify the response body
	expectedResponseBody := fmt.Sprintf(`{"access_token":"%s","message":"%s"}`, accessToken, jsonMessage)

	fmt.Println("expectedResponseBody : ", expectedResponseBody)

	if respRecorder.Body.String() != expectedResponseBody {
		t.Errorf("Handler returned unexpected body: got %v want %v", respRecorder.Body.String(), expectedResponseBody)
	}

}
