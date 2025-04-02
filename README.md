# Integración Continua con GitHub Actions y Webhook en Gin

Este proyecto implementa un servidor de Discord que se integra con GitHub Actions y webhooks para proporcionar notificaciones en tiempo real sobre eventos de desarrollo.

## Requisitos

- Go 1.21 o superior
- Token de bot de Discord
- ID del servidor de Discord (Guild ID)
- GitHub Actions configurado en tu repositorio

## Configuración

1. Crea un bot en Discord:
   - Ve a [Discord Developer Portal](https://discord.com/developers/applications)
   - Crea una nueva aplicación
   - Ve a la sección "Bot" y crea un bot
   - Copia el token del bot

2. Configura las variables de entorno:
   ```bash
   export DISCORD_TOKEN="tu-token-de-discord"
   export GUILD_ID="id-de-tu-servidor"
   ```

3. Configura los webhooks en GitHub:
   - Ve a la configuración de tu repositorio
   - En la sección "Webhooks", añade un nuevo webhook
   - URL: `http://tu-servidor:8080/webhook/github`
   - Content type: `application/json`
   - Selecciona los eventos:
     - Pull requests
     - Workflow runs

## Canales de Discord

El bot creará automáticamente tres canales:

1. **desarrollo**: Notificaciones de:
   - Pull requests (nuevos, reopened, ready_for_review)
   - Fusiones exitosas de PRs

2. **pruebas**: Notificaciones de:
   - Estado de los workflows de GitHub Actions
   - Resultados de las pruebas

3. **general**: Canal para comunicación general del equipo

## Ejecución

```bash
go run main.go
```

El servidor se iniciará en el puerto 8080.

## Integración con GitHub Actions

Para integrar con GitHub Actions, asegúrate de que tu workflow incluya los pasos necesarios para enviar notificaciones al webhook:

```yaml
name: CI/CD Pipeline

on:
  pull_request:
    types: [opened, reopened, ready_for_review]
  workflow_run:
    types: [completed]

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - name: Notify Discord
        run: |
          curl -X POST http://tu-servidor:8080/webhook/actions \
            -H "Content-Type: application/json" \
            -d '{"event": "${{ github.event_name }}", "status": "${{ job.status }}"}'
``` 