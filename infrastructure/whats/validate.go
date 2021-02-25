package whats

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

// ValidateNum verify if number is valid
func ValidateNum(number string, wac *whatsapp.Conn) (string, error) {
	suffix := "@c.us"
	payload := struct {
		Status int
		JID    string
	}{}

	resp, err := wac.Exist(fmt.Sprintf("%s%s", number, suffix))
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(<-resp), &payload); err != nil {
		return "", err
	}

	if payload.Status != 200 {
		return "", fmt.Errorf("status: %d", payload.Status)
	}

	return strings.TrimSuffix(payload.JID, suffix), nil
}
