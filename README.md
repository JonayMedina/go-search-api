# Music Search API

API de búsqueda de música que integra múltiples proveedores de servicios musicales.

## Descripción

Este proyecto es una API REST que permite buscar canciones a través de diferentes proveedores de servicios musicales, implementando caché y almacenamiento persistente de búsquedas.

## Características

- Búsqueda simultánea en múltiples proveedores
- Sistema de caché con Redis
- Almacenamiento persistente en MongoDB
- Autenticación JWT
- Middleware CORS
- Configuración basada en variables de entorno

## Tecnologías

- Go 1.23+
- MongoDB
- Redis
- Gin Web Framework
- JWT

## Estructura del Proyecto 

music-search-api/
├── cmd/
│ └── api/ # Punto de entrada de la aplicación
├── internal/
│ ├── config/ # Configuración de la aplicación
│ ├── cors/ # Middleware CORS
│ ├── errors/ # Manejo de errores
│ ├── handlers/ # Manejadores HTTP
│ ├── models/ # Modelos de datos
│ └── services/ # Lógica de negocio
└── go.mod

## Configuración

Variables de entorno requeridas:
env
SERVER_PORT=:8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=musicdb
REDIS_URI=localhost:6379
REDIS_PASSWORD=
JWT_SECRET=your-secret-key

## Instalación

1. Clonar el repositorio
bash
git clone `https://github.com/JonayMedina/music-search-api.git`

2. Instalar dependencias
bash
go mod download

3. Configurar variables de entorno

4. Ejecutar para iniciar la aplicación
bash
go run cmd/api/main.go

5. Acceder a la documentación de la API
La documentación se encuentra en el directorio: services-definition/services/

## Endpoints

### Públicos

- `GET /health` - Estado del servicio
- `POST /auth/login` - Inicio de sesión
- `POST /auth/register` - Registro de usuario

### Protegidos (requieren JWT)

- `GET /api/search` - Búsqueda de canciones
- `GET /api/history` - Historial de búsquedas

## Dependencias Principales

go
require (
github.com/gin-gonic/gin
github.com/redis/go-redis/v9
go.mongodb.org/mongo-driver/mongo
github.com/caarlos0/env/v6
)

## Desarrollo

Para ejecutar en modo desarrollo:

bash
go run cmd/api/main.go

## Tests

Ejecutar tests:

bash
go test -v ./...

## Contribución

1. Fork el proyecto
2. Crea tu rama de características (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.
