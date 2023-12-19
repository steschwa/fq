package firebase

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

var (
	ErrEmulatorNotRunning = errors.New("emulator not running")
)

func GetEmulatorHubHost() string {
	return os.Getenv(FirebaseEmulatorHubEnvName)
}

func IsRunningInEmulator() bool {
	emulatorHub := GetEmulatorHubHost()
	return emulatorHub != ""
}

func EmulatorDisableFunctionTriggers() error {
	return putEmulator("disableBackgroundTriggers")
}

func EmulatorEnableFunctionTriggers() error {
	return putEmulator("enableBackgroundTriggers")
}

func putEmulator(function string) error {
	if !IsRunningInEmulator() {
		return ErrEmulatorNotRunning
	}
	emulatorHub := GetEmulatorHubHost()

	url := fmt.Sprintf("http://%s/functions/%s", emulatorHub, function)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	client.Do(req)
	return nil
}
