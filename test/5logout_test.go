package test

import (
	"encoding/json"
	"example/auth/internal/pkg/authorize"
	"example/auth/internal/pkg/controller"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example/auth/test/dataset"

	"github.com/gin-gonic/gin"
)

func TestLogout(t *testing.T) {
	router := gin.Default()

	router.POST("/authorize/logout", authorize.ExtractJwtMiddleware(), controller.Logout)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/authorize/logout", nil)

	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	/* IMPORTANT : An active acces token should be replaced or an invalid one but is should have a valid refresh token */
	req.Header.Set("Authorization", "Bearer "+dataset.AccessToken)

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
	jsonMessage := "user logged out succesfully"

	// Verify the response body
	expectedResponseBody := fmt.Sprintf(`{"access_token":"%s","message":"%s"}`, accessToken, jsonMessage)

	fmt.Println("expectedResponseBody : ", expectedResponseBody)

	if respRecorder.Body.String() != expectedResponseBody {
		t.Errorf("Handler returned unexpected body: got %v want %v", respRecorder.Body.String(), expectedResponseBody)
	}

}
