package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ucho456job/passkey_sample/pkg/config"
	"github.com/ucho456job/passkey_sample/pkg/handler"
)

func init() {
	config.InitDB()
	config.InitRedis()
	config.InitWebAuthn()
}

func main() {
	router := gin.Default()

	router.Static("/static", "./template/build/static")
	router.GET("/", func(c *gin.Context) {
		c.File("./template/build/index.html")
	})

	router.POST("/api/auth/register-challenge", handler.ChallengeForRegister)
	router.POST("/api/auth/passkey", handler.Register)
	router.POST("/api/auth/login-challenge", handler.ChallengeForLogin)
	router.POST("/passkey-session", handler.Login)

	router.Run(":8080")
}
