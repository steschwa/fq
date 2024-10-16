package firestore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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

	if isEmulatorProject(projectID) {
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

func isEmulatorProject(projectID string) bool {
	// https://firebase.google.com/docs/emulator-suite/connect_firestore#choose_a_firebase_project
	return strings.HasPrefix(projectID, "demo-")
}

func setupEmulatorEnvironment() {
	emulatorHostEnv := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if emulatorHostEnv == "" {
		os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	}
}
