package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
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
	return u.Credentials
}

type UserData struct {
	ID   string `gorm:"primary_key;column:user_id"`
	Name string `gorm:"column:name"`
}

func ChallengeForRegister(c *gin.Context) {
	// Get user data
	userID := "1b2fa70d-8416-42c3-a789-96c4817129ea"
	var userData UserData
	if result := config.DB.Table("users").Where("user_id = ?", userID).First(&userData); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database for users"})
		return
	}
	user := WebAuthnUser{
		ID:          []byte(userID),
		Name:        userData.Name,
		DisplayName: userData.Name,
	}

	// Get existing credentials
	var existingCredentials []PublicKeyCredential
	if result := config.DB.Table("public_key_credentials").Where("user_id = ?", userID).Find(&existingCredentials); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database for existing credentials"})
		return
	}
	excludeCredentials := make([]protocol.CredentialDescriptor, len(existingCredentials))
	for i, cred := range existingCredentials {
		credentialID, err := base64.RawURLEncoding.DecodeString(cred.CredentialID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode credential ID"})
			return
		}
		excludeCredentials[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: credentialID,
		}
	}

	// Begin registration
	options, sessionData, err := config.WebAuthn.BeginRegistration(
		&user,
		webauthn.WithExclusions(excludeCredentials),
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationRequired,
		}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin registration"})
		return
	}

	// Save session data to Redis
	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal session data"})
		return
	}
	sessionKey := fmt.Sprintf("webauthn_challenge_register:%s", user.ID)
	err = config.Redis.Set(context.Background(), sessionKey, sessionDataJSON, 5*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session data to Redis"})
		return
	}

	c.JSON(http.StatusOK, options)
}

func LoginOption() webauthn.LoginOption {
	return func(opts *protocol.PublicKeyCredentialRequestOptions) {
		opts.UserVerification = protocol.VerificationRequired
	}
}

func ChallengeForLogin(c *gin.Context) {
	// Begin login
	options, sessionData, err := config.WebAuthn.BeginDiscoverableLogin(LoginOption())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin registration"})
		return
	}

	// Save session data to Redis
	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal session data"})
		return
	}
	sessionKey := fmt.Sprintf("webauthn_challenge_login:%s", options.Response.Challenge)
	err = config.Redis.Set(context.Background(), sessionKey, sessionDataJSON, 60*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session data to Redis"})
		return
	}

	c.JSON(http.StatusOK, options)
}
