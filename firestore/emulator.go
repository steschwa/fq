package firestore

import (
	"os"
	"strings"
)

func IsEmulatorProject(projectID string) bool {
	// https://firebase.google.com/docs/emulator-suite/connect_firestore#choose_a_firebase_project
	return strings.HasPrefix(projectID, "demo-")
}

func setupEmulatorEnvironment() {
	emulatorHostEnv := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if emulatorHostEnv == "" {
		os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	}
}
