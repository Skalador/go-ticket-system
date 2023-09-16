package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Skalador/go-ticket-system/models"
)

func setTestWorkingDirectory(t *testing.T) func() {
	t.Helper()

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Change the working directory to the root of the project
	err = os.Chdir("../")
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Return a cleanup function to reset the working directory
	return func() {
		err := os.Chdir(wd)
		if err != nil {
			t.Fatalf("Failed to reset working directory: %v", err)
		}
	}
}

func TestTicketRootHandler(t *testing.T) {
	// Create a mock Tickets object
	var ticketsCache models.Tickets
	models.InitTicketsCache(&ticketsCache)

	// Set the test working directory and defer cleanup
	resetWorkingDirectory := setTestWorkingDirectory(t)
	defer resetWorkingDirectory()

	// Create a request
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler function
	handler := TicketRootHandler(&ticketsCache)

	// Serve the request to the handler
	handler(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}
}
