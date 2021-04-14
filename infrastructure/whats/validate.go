package whats

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

// ValidateNum verify if number is valid
func ValidateNum(number string, wac *whatsapp.Conn) (string, error) {
	retsuffix := "@c.us"
	suffix := "@s.whatsapp.net"
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

	fmt.Printf("Payload from ValidateNum: [%+v]", payload)

	if payload.Status != 200 {
		return "", fmt.Errorf("status: %d", payload.Status)
	}

	return fmt.Sprintf("%s%s", strings.TrimSuffix(payload.JID, retsuffix), suffix), nil
	/*
		if len(num) > 0 {
			num2 := strings.Split(num, "@")
			destino2 = fmt.Sprintf("%s%s", num2[0], "@s.whatsapp.net")
		}


		return strings.TrimSuffix(payload.JID, suffix), nil
	*/
}
