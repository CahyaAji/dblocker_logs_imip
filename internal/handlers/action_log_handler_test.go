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

func setupActionLogHandler(t *testing.T) (*handlers.ActionLogHandler, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_action",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	repo := repository.NewActionLogRepository(gormDB)
	handler := handlers.NewActionLogHandler(repo)

	return handler, mock
}

func TestGetActionLogs(t *testing.T) {
	handler, mock := setupActionLogHandler(t)

	timestamp := time.Now()
	rows := sqlmock.NewRows([]string{"id", "timestamp", "user_id", "user_name", "action", "success", "detail"}).
		AddRow(1, timestamp, 101, "admin", "LOGIN", true, []byte("{}")).
		AddRow(2, timestamp, 102, "operator", "LOGOUT", true, []byte("{}"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "action_logs"`)).WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/action-logs", handler.GetActionLogs)

	req, _ := http.NewRequest(http.MethodGet, "/action-logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestCreateActionLog(t *testing.T) {
	handler, mock := setupActionLogHandler(t)

	input := models.ActionLog{
		UserID:   101,
		UserName: "admin",
		Action:   "REBOOT",
		// Success:  true,
		Detail: datatypes.JSON([]byte(`{"reason": "maintenance"}`)),
	}
	jsonBytes, _ := json.Marshal(input)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "action_logs"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	r := gin.Default()
	r.POST("/action-logs", handler.CreateActionLog)

	req, _ := http.NewRequest(http.MethodPost, "/action-logs", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Nil(t, mock.ExpectationsWereMet())
}
