package v1

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"database/sql"

	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/api/v1"
)

const (
	// apiVersion is the version of API provided by server
	apiVersion = "v1";
)

// toDoServiceServer is the implementation of v1.ToDoServiceServer proto interface
type toDOServiceServer struct {
	db *sql.DB
}

// NewToDoServiceServer creates ToDo Service
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDOServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
	// If API version is blank ("") then use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "Unsupported API version: service implements API version '%s' but asked for '%s' ", apiVersion, api)
		}
	}
	return nil
}

// connect returns a databse connection from pool
func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}