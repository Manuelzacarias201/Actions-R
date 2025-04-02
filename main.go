package main

import (
	"log"
	"os"

	"github.com/actions/internal/adapters/primary/http"
	"github.com/actions/internal/adapters/secondary/discord"
	"github.com/actions/internal/core/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Obtener variables de entorno
	discordToken := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("GUILD_ID")

	// Inicializar el adaptador de Discord
	discordAdapter, err := discord.NewDiscordAdapter(discordToken, guildID)
	if err != nil {
		log.Fatal("Error al crear el adaptador de Discord:", err)
	}

	// Inicializar Discord y crear canales
	if err := discordAdapter.Initialize(); err != nil {
		log.Fatal("Error al inicializar Discord:", err)
	}

	// Crear el servicio de notificaciones
	notificationService := services.NewNotificationService(discordAdapter, discordAdapter)

	// Crear el manejador de webhooks
	webhookHandler := http.NewWebhookHandler(notificationService)

	// Inicializar el servidor Gin
	r := gin.Default()

	// Configurar las rutas para los webhooks
	r.POST("/webhook/github", webhookHandler.HandleGitHubWebhook)
	r.POST("/webhook/actions", webhookHandler.HandleActionsWebhook)

	// Iniciar el servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
