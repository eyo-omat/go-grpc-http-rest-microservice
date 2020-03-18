package main

import (
	"time"
	"context"
	"log"
	"flag"
	
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	"github.com/eyo-omat/go-grpc-http-rest-microservice/pkg/api/v1"
)

const (
	// apiVersion is the version of the API provided by the server
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC servre in format host:port")
	flag.Parse()

	// set up a conncetion to the server
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	// Call to create
	createReq := v1.CreateRequest {
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title: "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			Reminder: reminder,
		},
	}

	createResp, err := c.Create(ctx, &createReq)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create Result: <%+v>\n\n", createResp)

	id := createResp.Id

	// Call to Read
	readReq := v1.ReadRequest {
		Api: apiVersion,
		Id: id,
	}

	readResp, err := c.Read(ctx, &readReq)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	log.Printf("Read Result: <%+v>\n\n", readResp)

	// Call to update
	updateReq := v1.UpdateRequest {
		Api: apiVersion,
		ToDo: &v1.ToDo {
			Id: readResp.ToDo.Id,
			Title: readResp.ToDo.Title,
			Description: readResp.ToDo.Description + "+ updated	",
			Reminder: readResp.ToDo.Reminder,
		},
	}

	updateResp, err := c.Update(ctx, &updateReq)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Update Result: <%+v>\n\n", updateResp)

	// Call to read all
	readAllReq := v1.ReadAllRequest {
		Api: apiVersion,
	}

	readAllResp, err := c.ReadAll(ctx, &readAllReq)
	if err != nil {
		log.Fatalf("Read all failed: %v", err)
	}
	log.Printf("ReadAll Result: <%+v>\n\n", readAllResp)

	// Call to delete
	deleteReq := v1.DeleteRequest {
		Api: apiVersion,
		Id: id,
	}

	deleteResp, err := c.Delete(ctx, &deleteReq)
	if err != nil {
		log.Fatalf("Delete failed %v", err)
	}
	log.Printf("Delete Result: <%+v>\n\n", deleteResp)


}