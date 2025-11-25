package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeonovDS/review-manager/internal/database"
	"github.com/LeonovDS/review-manager/internal/handlers"
	"github.com/LeonovDS/review-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	dbUser     string = "admin"
	dbPassword string = "admin"
	dbName     string = "test"
)

func runTest(t *testing.T, f func(t *testing.T, mux *http.ServeMux)) {
	ctx := t.Context()
	container, err := postgres.Run(ctx, "postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatal("Can't run container", "err", err)
		return
	}
	defer func() {
		_ = testcontainers.TerminateContainer(container)
	}()

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal("Can't generate connection string", "err", err)
		return
	}

	pool, err := database.Connect(ctx, connStr, "file://../migrations/")
	if err != nil {
		t.Fatal("Can't connect to database", "err", err)
		return
	}
	defer func() { pool.Close() }()
	router := handlers.NewRouter(pool)

	f(t, router)
}

// nolint:exhaustruct
func TestAddTeam(t *testing.T) {
	runTest(t, func(t *testing.T, mux *http.ServeMux) {
		team1 := model.Team{
			TeamName: "team1",
			Members: []model.User{
				{UserID: "u1", Username: "Alice", IsActive: true},
				{UserID: "u2", Username: "Bob", IsActive: true},
			},
		}

		team1Json, _ := json.Marshal(team1)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(team1Json))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		var resp model.Team
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		if err != nil {
			assert.NoError(t, err, "Should return correct json")
		}

		assert.Equal(t, rr.Code, http.StatusCreated)
		assert.Equal(t, team1, resp)
	})
}
