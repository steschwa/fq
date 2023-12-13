package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

func PrettifyJSON(jsonData []byte) (string, error) {
	indent := strings.Repeat(" ", 4)

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonData, "", indent); err != nil {
		return "", err
	}

	return prettyJSON.String(), nil
}
