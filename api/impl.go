package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// This is a simple example of how to make each of the api calls.
// A more realistic application might use a database to store and retrieve pets.
type Server struct {
	echo *echo.Echo
}

// NewServer returns a new Server.
func NewServer(e *echo.Echo) *Server {
	return &Server{echo: e}
}

// GetWorkload implements the GetWorkload interface.
func (s *Server) GetWorkload(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "GetWorkload")
}

// PostWorkload implements the PostWorkload interface.
func (s *Server) PostWorkload(ctx echo.Context) error {
	var workload Workload
	if err := ctx.Bind(&workload); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("error binding workload: %v", err))
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("PostWorkload: %v", workload))
}

// PutWorkload implements the PutWorkload interface.
func (s *Server) PutWorkload(ctx echo.Context) error {
	var workload Workload
	if err := ctx.Bind(&workload); err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("error binding workload: %v", err))
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("PutWorkload: %v", workload))
}

// DeleteWorkload implements the DeleteWorkload interface.
func (s *Server) DeleteWorkload(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "DeleteWorkload")
}
