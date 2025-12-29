package handlers_test

import (
	"bytes"
	"dblocker_logs_server/internal/handlers"
	"dblocker_logs_server/internal/models"
	"dblocker_logs_server/internal/repository"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupDeviceHandler creates a mock database and returns the handler and the mock object
func setupDeviceHandler(t *testing.T) (*handlers.DeviceHandler, sqlmock.Sqlmock) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Connect GORM to the mock DB
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	// Initialize the repository and handler with the mock DB
	repo := repository.NewDeviceRepository(gormDB)
	handler := handlers.NewDeviceHandler(repo)

	return handler, mock
}

func TestGetDevices(t *testing.T) {
	handler, mock := setupDeviceHandler(t)

	// Define the rows we expect the database to return
	rows := sqlmock.NewRows([]string{"id", "name", "type", "ip_address", "latitude", "longitude", "description"}).
		AddRow(1, "Living Room Cam", "camera", "192.168.1.50", -7.797, 110.370, "Main camera").
		AddRow(2, "Temp Sensor", "sensor", "192.168.1.51", -7.797, 110.370, "AC sensor")

	// Expect a SELECT query
	// We use regex because GORM might query specific columns or all columns (*)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "devices"`)).WillReturnRows(rows)

	// Setup Gin router in Test Mode
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/devices", handler.GetDevices)

	// Create a request
	req, _ := http.NewRequest(http.MethodGet, "/devices", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(data))

	// Ensure all expectations were met
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateDevice(t *testing.T) {
	handler, mock := setupDeviceHandler(t)

	input := models.Device{
		Name:      "New Sensor",
		Type:      "sensor",
		IPAddress: "192.168.1.99",
	}
	jsonBytes, _ := json.Marshal(input)

	// GORM Create usually starts a transaction
	mock.ExpectBegin()
	// Expect the INSERT query. Using regex to match broadly.
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "devices"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Setup Gin router
	r := gin.Default()
	r.POST("/devices", handler.CreateDevice)

	// Perform Request
	req, _ := http.NewRequest(http.MethodPost, "/devices", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}
