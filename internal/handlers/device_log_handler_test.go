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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDeviceLogHandler(t *testing.T) (*handlers.DeviceLogHandler, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_device_log",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	repo := repository.NewDeviceLogRepository(gormDB)
	handler := handlers.NewDeviceLogHandler(repo)

	return handler, mock
}

func TestGetDeviceLogs(t *testing.T) {
	handler, mock := setupDeviceLogHandler(t)

	timestamp := time.Now()
	rows := sqlmock.NewRows([]string{"id", "timestamp", "device_id", "device_name", "is_online", "status", "sensors"}).
		AddRow(1, timestamp, 5, "Sensor A", true, []byte("{}"), []byte("{}")).
		AddRow(2, timestamp, 6, "Sensor B", false, []byte("{}"), []byte("{}"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "device_logs"`)).WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/device-logs", handler.GetDeviceLogs)

	req, _ := http.NewRequest(http.MethodGet, "/device-logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateDeviceLog(t *testing.T) {
	handler, mock := setupDeviceLogHandler(t)

	input := models.DeviceLog{
		DeviceID:   5,
		DeviceName: "Sensor A",
		IsOnline:   true,
		Status:     datatypes.JSON([]byte(`{"battery": 90}`)),
		Sensors:    datatypes.JSON([]byte(`{"temp": 25}`)),
	}
	jsonBytes, _ := json.Marshal(input)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "device_logs"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	r := gin.Default()
	r.POST("/device-logs", handler.CreateDeviceLog)

	req, _ := http.NewRequest(http.MethodPost, "/device-logs", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}
