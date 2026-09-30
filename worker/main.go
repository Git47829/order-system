package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub/v2"
)

func main() {
	fmt.Println("Starting up Worker...")
	http.HandleFunc("/healthz", healthCheck)
	port := 8080
	fmt.Printf("Server running at http://localhost%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Errorf("Failed to start http Server %d", err)
	}

	var buf bytes.Buffer

	fmt.Println("Worker Ready")
	projectID := os.Getenv("projectID")
	subID := os.Getenv("subID")

	err = pullMsgs(&buf, projectID, subID)
	if err != nil {
		fmt.Errorf("failed to Process Message: %w", err)
	}
	processMsg(buf.String())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("healthy")
}

func processMsg(msg string) {
	r := rand.Intn(30)
	time.Sleep(time.Duration(r) * time.Second)
	fmt.Println("Processed Order: %w", msg)
}

func pullMsgs(w io.Writer, projectID, subID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	defer client.Close()

	sub := client.Subscriber(subID)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var recieved int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got Message: %q\n", string(msg.Data))
		atomic.AddInt32(&recieved, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Recieve %w", err)
	}
	fmt.Fprintf(w, "Recieved %d messages \n", recieved)

	return nil
}
