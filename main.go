package main
//p
import (
	"log"

	"github_webhook/api"
	"github_webhook/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Advertencia: No se pudo cargar el archivo .env: %v", err)
	}

	// Inicializar configuración
	cfg := config.NewConfig()

	// Crear router
	router := gin.Default()

	// Configurar middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Inicializar rutas
	api.SetupRoutes(router)

	// Iniciar servidor
	port := cfg.GetPort()
	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
//prueba
