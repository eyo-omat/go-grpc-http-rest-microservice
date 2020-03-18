package cmd

import (
	"database/sql"
	"fmt"
	"flag"
	"context"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/protocol/grpc"
	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/protocol/rest"
	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/service/v1"
)

// Config is configuration for the server
type Config struct {
	// gRPC Server start up parameters section
	// GRPCPort is TCP port to listen on by gRPC Server
	GRPCPort string

	// HTTP/REST start up parameters section
	// HTTPPort is the TCP port to listen on by HTTP/REST gateway
	HTTPPort string

	// DB Datastore parameters section
	// DatastoreDBHost is database host
	DatastoreDBHost string
	// DatastoreDBUser is the database username
	DatastoreDBUser string
	// DatastoreDBPassword is the database password
	DatastoreDBPassword string
	// DatastoreDBSchema is the database schema
	DatastoreDBSchema string
}

// RunServer runs gRPC server  and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get Configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	// Add MySQL driver specifc parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
			cfg.DatastoreDBUser,
			cfg.DatastoreDBPassword,
			cfg.DatastoreDBHost,
			cfg.DatastoreDBSchema,
			param)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}