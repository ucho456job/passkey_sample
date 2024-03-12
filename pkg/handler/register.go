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
	// ユーザー情報を取得
	userID := "dff8fd7b-a10f-4e33-8b60-a54d7ab4f5be"
	user := WebAuthnUser{
		ID:          []byte(userID),
		Name:        "test-email-01@example.com",
		DisplayName: "John Doe",
	}

	// セッションデータを取得
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

	// セッションデータを使って認証情報を取得
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

	if result := config.DB.Table("public_key_credentials").Create(&newCredential); result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Registration successful"})
}
