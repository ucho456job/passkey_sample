package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ucho456job/passkey_sample/pkg/config"
)

type WebAuthnUser struct {
	ID          []byte
	Name        string
	DisplayName string
	Icon        string
	Credentials []webauthn.Credential
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.ID
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.Name
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return []webauthn.Credential{}
}

func Challenge(c *gin.Context) {
	user := WebAuthnUser{
		ID:          []byte("dff8fd7b-a10f-4e33-8b60-a54d7ab4f5be"),
		Name:        "John Doe",
		DisplayName: "John Doe",
	}

	options, sessionData, err := config.WebAuthn.BeginRegistration(
		&user,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin registration"})
		return
	}

	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal session data"})
		return
	}
	sessionKey := fmt.Sprintf("webauthn_challenge_register:%s", user.ID)
	err = config.Redis.Set(context.Background(), sessionKey, sessionDataJSON, 0).Err() // 期限を設定しない場合は0
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session data to Redis"})
		return
	}

	c.JSON(http.StatusOK, options)
}
