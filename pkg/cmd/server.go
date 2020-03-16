package cmd

import (
	"database/sql"
	"fmt"
	"flag"
	"context"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/protocol/grpc"
	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/service/v1"
)

// Config is configuration for the server
type Config struct {
	// gRPC Server start parameters section
	// gRPC is TCP port to listen by gRPC Server
	GRPCPort string

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
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
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

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}