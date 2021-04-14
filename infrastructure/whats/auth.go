package whats

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/workstash/whapi/config"
	"github.com/workstash/whapi/helper/logger"

	"github.com/Rhymen/go-whatsapp"
)

//NewConn connects to whatsapp
func NewConn() (*whatsapp.Conn, error) {
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("creating connection: %w", err)
	}
	client := config.Main.Client
	wac.SetClientName(client.LongName, client.ShortName, client.Version)
	wac.SetClientVersion(2, 2021, 4)

	return wac, nil
}

//Auth try to connect to existing session
func auth(wac *whatsapp.Conn, sessionPath string) error {
	//load saved session
	session, err := readSession(sessionPath)
	if err != nil {
		return err
	}
	//restore session
	session, err = wac.RestoreWithSession(session)
	if err != nil {
		return err
	}
	return nil
}

func Auth(wac *whatsapp.Conn, sessionPath string) error {
	//load saved session
	session, err := readSession(sessionPath)
	if err != nil {
		return err
	}
	//restore session
	session, err = wac.RestoreWithSession(session)
	if err != nil {
		return err
	}
	return nil
}

//Login login to whats app account
func Login(wac *whatsapp.Conn, device string) (string, error) {
	//no saved session -> regular login
	sessionPath := fmt.Sprintf("%s/%s.gob", config.Main.API.SessionPath, device)
	qrCode := make(chan string, 10)
	if err := auth(wac, sessionPath); err != nil {
		go func() {
			//login
			session, err := wac.Login(qrCode)
			if err != nil {
				// qr code scan timed out
				return
			}
			//save session if login was successful
			if err := writeSession(session, sessionPath); err != nil {
				logger.Printf("writing session: %s", err)
			}
		}()
	} else {
		return "", nil
	}
	return <-qrCode, nil
}

func readSession(sessionPath string) (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(sessionPath)
	if err != nil {
		return session, err
	}
	defer file.Close()
	return session, gob.NewDecoder(file).Decode(&session)
}

func writeSession(session whatsapp.Session, sessionPath string) error {
	file, err := os.Create(sessionPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return gob.NewEncoder(file).Encode(session)
}
