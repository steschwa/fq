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

func setFirebaseEmulatorEnv() {
	envMap := map[string]string{
		"FIREBASE_EMULATOR_HUB":   "localhost:4000",
		"FIRESTORE_EMULATOR_HOST": "localhost:8080",
	}

	for key, val := range envMap {
		if envVal := os.Getenv(key); envVal == "" {
			os.Setenv(key, val)
		}
	}
}

func getIsEmulatorProject(projectID string) bool {
	return strings.HasPrefix(projectID, "demo-")
}
