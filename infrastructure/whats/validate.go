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

	//sessionPath := fmt.Sprintf("%s/%s.gob", config.Main.API.SessionPath, device)
	//if err := auth(wac, sessionPath); err == nil {

	fullnumber := fmt.Sprintf("%s%s", number, suffix)

	resp, err := wac.Exist(fullnumber)
	if err != nil {
		fmt.Println("Error calling wac.Exist: ", err.Error())
		return "", err
	}

	if err := json.Unmarshal([]byte(<-resp), &payload); err != nil {
		fmt.Println("Error unmarshalling Exist: ", err.Error())
		return "", err
	}

	fmt.Printf("Payload from ValidateNum: [%+v]", payload)

	ret := fmt.Sprintf("%s%s", strings.TrimSuffix(payload.JID, retsuffix), suffix)

	if payload.Status != 200 {
		fmt.Println("Payload not 200: ")
		fmt.Println("---------------------------------------------------------------------------")
		fmt.Println("ValidateNum")
		fmt.Println("Full Number:", fullnumber)
		fmt.Printf("Payload: [%+v]\n", payload)
		fmt.Println("retorno:", ret)
		fmt.Println("---------------------------------------------------------------------------")
		return "", fmt.Errorf("status: %d", payload.Status)
	}

	return ret, nil
	//}
	/*
		if len(num) > 0 {
			num2 := strings.Split(num, "@")
			destino2 = fmt.Sprintf("%s%s", num2[0], "@s.whatsapp.net")
		}


		return strings.TrimSuffix(payload.JID, suffix), nil
	*/
}
