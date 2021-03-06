package v1

import (
	"time"
	"fmt"
	"context"
	"database/sql"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/api/v1"
)

const (
	// apiVersion is the version of API provided by server
	apiVersion = "v1";
)

// toDoServiceServer is the implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {
	db *sql.DB
}

// NewToDoServiceServer creates ToDo Service
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db: db}
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

// Create a new task
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// Validate requested API version is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get database connection
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format->"+err.Error())
	}

	// insert ToDo entity data
	res, err := c.ExecContext(ctx, "INSERT INTO ToDo(`Title`, `Description`, `Reminder`) VALUES(?,?,?)", req.ToDo.Title, req.ToDo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into ToDO-> "+err.Error())
	}

	// get ID of created Task
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo -> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id: id,
	}, nil
}

// Read a task
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// Validate requested API version is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get database connection
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// Retrieve Todo by ID
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo -> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo->"+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.Id))
	}

	// Retrieve ToDo Data
	var td v1.ToDo
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+ err.Error())
	}
	td.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format->"+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple ToDo rows with ID='%d'", req.Id))
	}

	return &v1.ReadResponse {
		Api: apiVersion,
		ToDo: &td,
	}, nil
}

// Update a task
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// Validate requested API version is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get database connection
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format->"+err.Error())
	}

	// update todo
	res, err := c.ExecContext(ctx, "UPDATE ToDo SET `Title`=?, `Description`=?, `Reminder`=? WHERE `ID`=?", req.ToDo.Title, req.ToDo.Description, reminder, req.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo->"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.ToDo.Id))
	}

	return &v1.UpdateResponse {
		Api: apiVersion,
		Updated: rows,
	}, nil
}

// Delete a task
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// Validate requested API version is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get database connection
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete todo task
	res, err := c.ExecContext(ctx, "DELETE FROM ToDo WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete Todo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.Id))
	}

	return &v1.DeleteResponse {
		Api: apiVersion,
		Deleted: rows,
	}, nil
}

// Read all todo tasks
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// Validate requested API version is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get database connection
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get all todos as a list
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	var reminder time.Time
	list := []*v1.ToDo{}
	for rows.Next(){
		td := new (v1.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field value for Todo row-> "+ err.Error())
		}
		td.Reminder, err = ptypes.TimestampProto(reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has an invalid format-> "+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+ err.Error())
	}

	return &v1.ReadAllResponse {
		Api: apiVersion,
		ToDos: list,
	}, nil
}