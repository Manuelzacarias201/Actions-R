package api
//ss
import (
	"github_webhook/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	webhookGroup := router.Group("/webhook")
	{
		webhookGroup.POST("/github", services.HandleGitHubWebhook)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Servicio funcionando correctamente",
	})
}
