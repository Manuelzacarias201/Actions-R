package api

//ss
import (
	"github_webhook/services"
	"github_webhook/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Ruta de prueba
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "El servidor está funcionando correctamente",
		})
	})

	discordService := utils.NewDiscordService()
	handler := services.NewWebhookHandler(discordService)

	webhookGroup := router.Group("/webhook")
	{
		webhookGroup.POST("/github", handler.HandleGitHubWebhook)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Servicio funcionando correctamente",
	})
}
