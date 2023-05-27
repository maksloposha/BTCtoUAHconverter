package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAddUser(t *testing.T) {
	// Prepare a test user
	user := User{
		ID:    1,
		Email: "test@example.com",
	}
	userJSON, _ := json.Marshal(user)

	// Create a request with the test user JSON
	request, _ := http.NewRequest("POST", "/subscribe", bytes.NewBuffer(userJSON))

	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the addUser function with the recorder and request
	addUser(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, recorder.Code)
	}

	// Check the response body
	expectedResponse := "User added successfully"
	if recorder.Body.String() != expectedResponse {
		t.Errorf("Expected response body '%s', but got '%s'", expectedResponse, recorder.Body.String())
	}

	// Check if the user exists in the file
	exists, _ := userExists(user)
	if !exists {
		t.Error("Expected user to exist in the file, but it doesn't")
	}

	// Clean up the test file
	os.Remove("users.json")
}

func TestUserExists(t *testing.T) {
	// Prepare a test user
	user := User{
		ID:    1,
		Email: "test@example.com",
	}

	// Add the user to the file
	_ = addUserToFile(user)

	// Check if the user exists
	exists, err := userExists(user)
	if err != nil {
		t.Errorf("Error checking if user exists: %s", err.Error())
	}

	if !exists {
		t.Error("Expected user to exist, but it doesn't")
	}

	// Clean up the test file
	os.Remove("users.json")
}

func addUserToFile(user User) error {
	// Open the file in append mode
	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert user to JSON
	jsonData, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return err
	}

	// Write user JSON to file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	// Write the closing square bracket
	_, err = file.WriteString("\n]\n")
	if err != nil {
		return err
	}

	return nil
}
