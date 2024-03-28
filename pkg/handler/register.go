package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ucho456job/passkey_sample/pkg/config"
)

type PublicKeyCredential struct {
	CredentialID    string `gorm:"primary_key;column:credential_id"`
	UserID          string `gorm:"column:user_id"`
	PublicKey       string `gorm:"column:public_key"`
	AttestationType string `gorm:"column:attestation_type"`
	SignCount       uint32 `gorm:"column:sign_count"`
	AAGUID          string `gorm:"column:aagu_id"`
	Platform        string `gorm:"column:platform"`
	UserAgent       string `gorm:"column:user_agent"`
}

func Register(c *gin.Context) {
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

	// Get session data
	sessionKey := fmt.Sprintf("webauthn_challenge_register:%s", userID)
	result, err := config.Redis.Get(context.Background(), sessionKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session data from Redis"})
		return
	}
	var sessionData webauthn.SessionData
	err = json.Unmarshal([]byte(result), &sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal session data"})
		return
	}

	// Finish registration
	credential, err := config.WebAuthn.FinishRegistration(&user, sessionData, c.Request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish registration"})
		return
	}

	newCredential := PublicKeyCredential{
		CredentialID:    base64.RawURLEncoding.EncodeToString(credential.ID),
		UserID:          userID,
		PublicKey:       base64.RawURLEncoding.EncodeToString(credential.PublicKey),
		AttestationType: credential.AttestationType,
		SignCount:       credential.Authenticator.SignCount,
		AAGUID:          base64.StdEncoding.EncodeToString(credential.Authenticator.AAGUID),
		Platform:        "platform_dummy",
		UserAgent:       "user_agent_dummy",
	}

	// Save credential to database
	if result := config.DB.Table("public_key_credentials").Create(&newCredential); result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Registration successful"})
}
