package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type (
	JSONSerializable interface {
		ToJSON() (string, error)
	}
)

func ToJSON(data any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return "", errors.New("failed to serialize firestore docs to json")
	}

	if IsInteractiveTTY() {
		jsonString, err := prettifyJSON(jsonData)
		if err != nil {
			jsonString = string(jsonData)
		}

		return jsonString, nil
	}

	return string(jsonData), nil
}

func prettifyJSON(jsonData []byte) (string, error) {
	indent := strings.Repeat(" ", 4)

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonData, "", indent); err != nil {
		return "", err
	}

	return prettyJSON.String(), nil
}
