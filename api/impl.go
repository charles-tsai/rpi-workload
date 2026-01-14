//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml -o api.gen.go api.yaml
package api

import (
	"context"
)

// Server implements the StrictServerInterface generated from the OpenAPI spec.
type Server struct{}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

// GetApps implements the GetApps interface.
func (s *Server) GetApps(ctx context.Context, request GetAppsRequestObject) (GetAppsResponseObject, error) {
	// For now, return a fixed list of apps
	return GetApps200JSONResponse{
		{Id: "1", Name: "App 1"},
		{Id: "2", Name: "App 2"},
	}, nil
}

// CreateApp implements the CreateApp interface.
func (s *Server) CreateApp(ctx context.Context, request CreateAppRequestObject) (CreateAppResponseObject, error) {
	// Echo back the created app
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
