//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml -o api.gen.go api.yaml
package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBInterface defines the methods we need from pgxpool.Pool
type DBInterface interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
}

// Server implements the StrictServerInterface generated from the OpenAPI spec.
type Server struct {
	db DBInterface
}

// NewServer returns a new Server.
func NewServer(db DBInterface) *Server {
	return &Server{db: db}
}

// InitSchema creates the necessary tables if they don't exist
func (s *Server) InitSchema(ctx context.Context) error {
	_, err := s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS apps (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create apps table: %w", err)
	}
	return nil
}

// GetApps implements the GetApps interface.
func (s *Server) GetApps(ctx context.Context, request GetAppsRequestObject) (GetAppsResponseObject, error) {
	rows, err := s.db.Query(ctx, "SELECT id, name FROM apps")
	if err != nil {
		return nil, fmt.Errorf("failed to query apps: %w", err)
	}
	defer rows.Close()

	var apps []App
	for rows.Next() {
		var app App
		if err := rows.Scan(&app.Id, &app.Name); err != nil {
			return nil, fmt.Errorf("failed to scan app: %w", err)
		}
		apps = append(apps, app)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return GetApps200JSONResponse(apps), nil
}

// CreateApp implements the CreateApp interface.
func (s *Server) CreateApp(ctx context.Context, request CreateAppRequestObject) (CreateAppResponseObject, error) {
	_, err := s.db.Exec(ctx, "INSERT INTO apps (id, name) VALUES ($1, $2)", request.Body.Id, request.Body.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to insert app: %w", err)
	}

	return CreateApp201JSONResponse{
		Id:   request.Body.Id,
		Name: request.Body.Name,
	}, nil
}

// GetAppById implements the GetAppById interface.
func (s *Server) GetAppById(ctx context.Context, request GetAppByIdRequestObject) (GetAppByIdResponseObject, error) {
	// Return a dummy app
	return GetAppById200JSONResponse{
		Id:   request.AppId,
		Name: "Dummy App",
	}, nil
}

// UpdateApp implements the UpdateApp interface.
func (s *Server) UpdateApp(ctx context.Context, request UpdateAppRequestObject) (UpdateAppResponseObject, error) {
	// Echo back the updated app
	return UpdateApp200JSONResponse{
		Id:   request.AppId,
		Name: request.Body.Name,
	}, nil
}
