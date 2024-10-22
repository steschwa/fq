package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

const (
	timeoutNewClient = 10
)

func NewClient(projectID string) (*firestore.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeoutNewClient)
	defer cancel()

	if IsEmulatorProject(projectID) {
		setupEmulatorEnvironment()
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithTelemetryDisabled())
	if errors.Is(err, context.Canceled) {
		return nil, fmt.Errorf("timed-out after %d seconds", timeoutNewClient)
	}
	if err != nil {
		return nil, fmt.Errorf("creating firestore client")
	}

	return client, nil
}
