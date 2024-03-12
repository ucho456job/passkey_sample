package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/ucho456job/passkey_sample/pkg/config"
)

func Login(c *gin.Context) {
	parsedResponse, err := protocol.ParseCredentialRequestResponse(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse credential request response"})
		return
	}

	sessionKey := fmt.Sprintf("webauthn_challenge_login:%s", parsedResponse.Response.CollectedClientData.Challenge)
	result, err := config.Redis.Get(c, sessionKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session data"})
		return
	}
	var sessionData webauthn.SessionData
	err = json.Unmarshal([]byte(result), &sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal session data"})
		return
	}

	handler := func(rawID, userHandle []byte) (user webauthn.User, err error) {
		var creds []PublicKeyCredential
		userID := string(userHandle)
		if result := config.DB.Table("public_key_credentials").Where("user_id = ?", userID).Find(&creds); result.Error != nil {
			return nil, fmt.Errorf("failed to query database for credentials: %w", result.Error)
		}

		webAuthnCreds := make([]webauthn.Credential, len(creds))
		for i, c := range creds {
			publicKey, err := base64.RawURLEncoding.DecodeString(c.PublicKey)
			if err != nil {
				return nil, fmt.Errorf("failed to decode public key: %w", err)
			}
			credentialID, err := base64.RawURLEncoding.DecodeString(c.CredentialID)
			if err != nil {
				return nil, fmt.Errorf("failed to decode credential ID: %w", err)
			}

			webAuthnCreds[i] = webauthn.Credential{
				ID:              credentialID,
				PublicKey:       publicKey,
				AttestationType: c.AttestationType,
				Authenticator: webauthn.Authenticator{
					SignCount: c.SignCount,
				},
			}
		}

		user = &WebAuthnUser{
			ID:          userHandle,
			Name:        "John Doe",
			DisplayName: "John Doe",
			Icon:        "",
			Credentials: webAuthnCreds,
		}
		fmt.Println("userID: ", userID)
		return user, nil
	}
	_, err = config.WebAuthn.ValidateDiscoverableLogin(handler, sessionData, parsedResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
}
