package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func waitEnvoy(url string, hostHeader string) {
	client := &http.Client{
		Timeout: 100 * time.Millisecond,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create a HTTP client: %v", err)
		return // Something bad happened. Skip waiting.
	}
	req.Host = hostHeader

	for {
		err := sendReq(client, req)
		if err == nil {
			return // Succeed.
		}

		log.Printf("Failed to wait envoy, retring: %v", err)
		time.Sleep(300 * time.Millisecond)
	}
}

func sendReq(client *http.Client, req *http.Request) error {
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send the HTTP request: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Println("Got 200 from envoy")
		// XXX: Read response body to close the socket properly.
		_, err := io.Copy(ioutil.Discard, res.Body)
		if err != nil {
			log.Printf("Failed to read the response body: %v", err)
			return nil // Anyway we succeeded to check Envoy's availability.
		}
		return nil // Succeed.
	}

	return fmt.Errorf("Status code is not 200: %v", res.StatusCode)
}

func execCmd(args []string) {
	if len(args) < 1 {
		log.Fatalf("At least one command line argument is required")
	}
	path := args[0]
	cmd, err := exec.LookPath(path)
	if err != nil {
		log.Fatalf("Command not found: %v", path)
	}

	log.Printf("Call execve with: %v", args)
	// Does not return if succeeds
	e := syscall.Exec(cmd, args, os.Environ())
	if e != nil {
		log.Fatalf("%v", e)
	}
}

func main() {
	// Timeout for overall operations.
	timeoutPrt := flag.Int("timeout", 10000, "timeout msec")
	healthcheckURLPtr := flag.String("envoy-healthcheck-url", "", "Healthcheck URL to access upstream service via Envoy (required)")
	hostHeaderPtr := flag.String("envoy-host-header", "", "HTTP Host header which represents an upstream service (required)")
	flag.Parse()
	if *healthcheckURLPtr == "" {
		log.Fatalf("envoy-healthcheck-url flag is required")
	}
	if *hostHeaderPtr == "" {
		log.Fatalf("envoy-host-header flag is required")
	}
	timeout := time.Millisecond * time.Duration(*timeoutPrt)

	c := make(chan int, 1)
	go func() {
		waitEnvoy(*healthcheckURLPtr, *hostHeaderPtr)
		c <- 0
	}()
	select {
	case <-c:
		log.Printf("Waiting envoy succeeded")
	case <-time.After(timeout):
		log.Printf("Waiting envoy timed out")
	}

	execCmd(flag.Args())
}
