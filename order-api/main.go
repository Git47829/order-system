package main

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
)

func main() {
}

func publish(w io.Writer, projectID, topicID, msg string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	publisher := client.Publisher(topicID)
	result := publisher.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})

	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}
