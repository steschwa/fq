package firebase

import (
	"os"
	"strings"
)

func setupFirebase(projectID string) {
	if getIsEmulatorProject(projectID) {
		setFirebaseEmulatorEnv()
	}
}

const (
	FirebaseEmulatorHubEnvName   = "FIREBASE_EMULATOR_HUB"
	FirestoreEmulatorHostEnvName = "FIRESTORE_EMULATOR_HOST"
)

func setFirebaseEmulatorEnv() {
	envMap := map[string]string{}
	envMap[FirebaseEmulatorHubEnvName] = "localhost:4400"
	envMap[FirestoreEmulatorHostEnvName] = "localhost:8080"

	for key, val := range envMap {
		if envVal := os.Getenv(key); envVal == "" {
			os.Setenv(key, val)
		}
	}
}

func getIsEmulatorProject(projectID string) bool {
	return strings.HasPrefix(projectID, "demo-")
}
