package whats

import (
	"fmt"
	"github.com/workstash/whapi/config"

	"github.com/Rhymen/go-whatsapp"
)

//SendMessageA send message to the specified num with auth
func SendMessageA(wac *whatsapp.Conn, device, num, msg string) error {
	sessionPath := fmt.Sprintf("%s/%s.gob", config.Main.API.SessionPath, device)
	if err := auth(wac, sessionPath); err == nil {
		text := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: num + "@s.whatsapp.net",
			},
			Text: msg,
		}

		_, err = wac.Send(text)
		if err != nil && err.Error() != "sending message timed out" {
			return err
		}
	} else {
		return ErrConnecting
	}
	return nil
}

//SendMessage send message to the specified num without auth
func SendMessage(wac *whatsapp.Conn, num, msg string) error {
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: num + "@s.whatsapp.net",
		},
		Text: msg,
	}

	_, err := wac.Send(text)
	if err != nil && err.Error() != "sending message timed out" {
		return err
	}

	return nil
}
