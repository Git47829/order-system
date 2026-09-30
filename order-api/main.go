package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub/v2"
)

type Order struct {
	Item string `json:"order"`
}

func main() {
	http.HandleFunc("order/{order}", createOrder)
	http.HandleFunc("/healthz", healthCheck)

	port := 3000
	fmt.Printf("Server running at http://localhost%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Errorf("Failed to start http server %d", err)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	projectID := os.Getenv("projectID")
	topicID := os.Getenv("topicID")
	w.Header().Set("Content-Type", "application/json")
	var userOrder Order
	userOrder.Item = r.PathValue("order")

	if userOrder.Item == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Item Order cannot be empty",
		})
		return
	}

	publish(projectID, topicID, userOrder)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(userOrder)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("healthy")
}

func publish(projectID, topicID string, msg Order) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	publisher := client.Publisher(topicID)
	result := publisher.Publish(ctx, &pubsub.Message{
		Data: []byte(msg.Item),
	})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	return nil
}
