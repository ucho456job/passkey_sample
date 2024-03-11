package config

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

var WebAuthn *webauthn.WebAuthn

func InitWebAuthn() {
	var err error
	WebAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "test-rp",
		RPID:          "localhost",
		RPOrigin:      "http://localhost:8080",
	})
	if err != nil {
		panic(err)
	}
}
