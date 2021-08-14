package model

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

type Campaign struct {
	ApiKey      string `json:"ApiKey"`
	AdminApiKey string `json:"AdminApiKey"`
	EncKey      []byte `json:"EncKey"`
	WithZip     bool   `json:"WithZip"`
	WithEnc     bool   `json:"WithEnc"`

	ServerUrl string `json:"ServerUrl"` // URL of the server, as viewed from the clients

	PacketSendPath   string `json:"PacketSendPath"`
	PacketGetPath    string `json:"PacketGetPath"`
	FileUploadPath   string `json:"FileUploadPath"`
	FileDownloadPath string `json:"FileDownloadPath"`

	WebuiPath          string
	ClientUseWebsocket bool `json:"ClientUseWebsocket"`
}

func MakeCampaign() Campaign {
	apiKey := "Secret-Api-Key"
	adminApiKey := "Secret-AdminApi-Key"
	encKey := "Secret-Enc-Key"
	serverUrl := "http://localhost:8080"

	// Generate the actual AES key based on encKey
	key, err := scrypt.Key([]byte(encKey), []byte("antnium-salt"), 2048, 8, 1, 32)
	if err != nil {
		log.Fatal("Could not create key")
	}

	c := Campaign{
		apiKey,
		adminApiKey,
		key,
		true,
		true,
		serverUrl,
		"/send",
		"/get/",
		"/upload/",
		"/static/",
		"/webui",
		true,
	}
	return c
}
