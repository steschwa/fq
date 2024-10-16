package firestore

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	if isEmulatorProject(projectID) {
		setupEmulatorEnvironment()
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithTelemetryDisabled())
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
