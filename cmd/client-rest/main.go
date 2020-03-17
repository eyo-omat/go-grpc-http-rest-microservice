package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"strings"
	"net/http"
	"time"
	"flag"

)

func main() {
	// get configuration
	address := flag.String("server", "http://localhost:8080", "Http gateway url, eg http://localhost:8080")
	flag.Parse()

	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	var body string

	// Call create task
	resp, err := http.Post(*address+"/v1/todo", "application/json", strings.NewReader(
		fmt.Sprintf(`{
			"api":"v1",
			"toDo": {
				"title":"title (%s),
				"description": "description (%s),
				"reminder": "%s"
			}
		}`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to call Create method: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read create task response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// parse ID of created ToDo
	var created struct {
		API string `json:"api"`
		ID string `json:"id"`
	}
	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("failed to Unmarshal JSON response of create method: %v", err)
		fmt.Println("error:", err)
	}

	// Call Read Task
	readResp, err := http.Get(fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(readResp.Body)
	readResp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read Read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", readResp.StatusCode, body)

	// Call update task
	updateReq, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID), strings.NewReader(
		fmt.Sprintf(`{
			"api":"v1",
			"toDO": {
				"title":"title (%s),
				"description": "description (%s),
				"reminder": "%s"
			}
		}`, pfx, pfx, pfx)))
	updateReq.Header.Set("Content-Type", "application/json")
	updateresp, err := http.DefaultClient.Do(updateReq)
	if err != nil {
		log.Fatalf("failed to call update method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(updateresp.Body)
	updateresp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read Update response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Update response: Code=%d, Body=%s\n\n", updateresp.StatusCode, body)

	// Call Read all tasks
	readAllResp, err := http.Get(*address + "/v1/todo/all")
	if err != nil {
		log.Fatalf("failed to call ReadAll method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(readAllResp.Body)
	readAllResp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read ReadAll response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", readAllResp.StatusCode, body)

	// Call delete task
	deleteReq, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID), nil)
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		log.Fatalf("failed to call Delete method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(deleteResp.Body)
	deleteResp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read Delete response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Delete Response: Code=%d, Body=%s\n\n", updateresp.StatusCode, body)
}