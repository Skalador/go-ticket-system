package db

import (
	"net/http"
	"os"
	"testing"

	"github.com/Skalador/go-ticket-system/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
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

// TestAddTicketToCacheAndDB tests the AddTicketToCacheAndDB function.
func TestAddTicketToCacheAndDB(t *testing.T) {
	// Create a mock MongoDB client
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	// Set the test working directory and defer cleanup
	resetWorkingDirectory := setTestWorkingDirectory(t)
	defer resetWorkingDirectory()

	// Create a mock TicketsCache
	var ticketsCache models.Tickets
	models.InitTicketsCache(&ticketsCache)

	// Assert the expected behavior
	mt.Run("success", func(mt *mtest.T) {
		// Create a test HTTP request
		req := &http.Request{
			Form: map[string][]string{
				"subject":     {"Test Subject"},
				"description": {"Test Description"},
				"severity":    {"SEV3"},
			},
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		AddTicketToCacheAndDB(req, mt.Client, &ticketsCache)
		testTicket := models.Ticket{Subject: "Test Subject", Description: "Test Description", Severity: "SEV3", ID: findMaxID(&ticketsCache)}
		assert.Equal(t, testTicket, ticketsCache.Tickets[len(ticketsCache.Tickets)-1], "Entry in ticketsCache matches the desired result.")
		assert.Equal(t, 4, len(ticketsCache.Tickets), "Dummy data from data.json contains 3 entries. After adding one, the desired length is 4.")
	})
}

func TestDeleteTicketFromCacheAndDB(t *testing.T) {
	// Create a mock MongoDB client
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	// Set the test working directory and defer cleanup
	resetWorkingDirectory := setTestWorkingDirectory(t)
	defer resetWorkingDirectory()

	// Create a mock TicketsCache
	var ticketsCache models.Tickets
	models.InitTicketsCache(&ticketsCache)

	// Assert the expected behavior
	mt.Run("success", func(mt *mtest.T) {
		// Create a test HTTP request
		req := &http.Request{
			Form: map[string][]string{
				"subject":     {"Improve Ticket List displaying option"},
				"description": {"The ID and Severity should be highlighted in the ticket list."},
				"id":          {"1"},
				"severity":    {"SEV4"},
			},
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		DeleteTicketFromCacheAndDB(req, mt.Client, &ticketsCache)
		assert.Equal(t, 2, len(ticketsCache.Tickets), "Dummy data from data.json contains 3 entries. After removing one, the desired length is 2.")
	})
}
