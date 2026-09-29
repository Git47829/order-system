package worker

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub/v2"
)

func main() {
	fmt.Println("Starting up Worker...")
	fmt.Println("Worker Ready")
	//impliment logic for getting the project ID and SubID Dynamically

	msg, err := pullMsgs(w io.Writer, projectID string, subID string)
	processMsg(msg)
}

func processMsg(msg string) {
	r := rand.Intn(30)
	time.Sleep(time.Duration(r) * time.Second)
	fmt.Println("Processed Order: %w", msg)
}

func pullMsgs(w io.Writer, projectID, subID string) (string, error) {
	ch := make(chan string)
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return "Error connecting to PubSub", fmt.Errorf("pubsub.NewClient: %w", err)
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
		ch <- string(msg.Data)
	})
	if err != nil {
		return "Error Recieving Message", fmt.Errorf("sub.Recieve %w", err)
	}
	fmt.Fprintf(w, "Recieved %d messages \n", recieved)

	v := <-ch
	return v, nil
}
