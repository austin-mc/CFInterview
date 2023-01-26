package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

//const URL = "http://localhost:8080"

type Health struct {
	NumRequests int `json:"numRequests"`
	NumErrors   int `json:"numErrors"`
}

func main() {

}

func healthCheck(url string, maxTime int) error {

	c := &http.Client{
		Timeout: time.Duration(time.Duration(maxTime).Seconds()),
	}

	resp, err := c.Get(url + "/health")
	if err != nil {
		log.Println(err)
		return err
	}

	if os.IsTimeout(err) {
		return errors.New("request timeout")
	}

	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return errors.New("invalid request")
	}
	if statusCode >= 500 {
		return errors.New("server error")
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("response status code %d", statusCode)
	}

	body := resp.Body
	healthData := Health{}
	err = json.NewDecoder(body).Decode(&healthData)
	if err != nil {
		return err
	}

	return nil
}

func healthCheckMonitor(url string, cadence int, timeout int) {
	for {
		err := healthCheck(url, timeout)
		if err != nil {
			log.Printf("Health check returned an error: %s", err)
		}
		time.Sleep(time.Duration(time.Duration(cadence).Seconds()))
	}
}

// HTTP req to endpoint health check
// Returns if service is in a healthy state
// /health endpoint
