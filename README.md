# Greenhouse BE 🏡

Backend del sistema **Digital Greenhouse**, construido en **Go** utilizando una arquitectura limpia (Clean Architecture/ Onion Architecture) con enrutamiento ágil y persistencia relacional.

## 🚀 Tecnologías Principales
* **Lenguaje:** Go 1.25+
* **Enrutador HTTP:** Go-Chi (v5)
* **ORM:** GORM (MySQL/MariaDB driver)
* **Autenticación:** JWT (JSON Web Tokens)
* **Base de Datos:** MariaDB 10.11

---

## 🛠️ Configuración y Desarrollo

### Requisitos Previos
* Go instalado (versión 1.25 o superior)
* Docker y Docker Compose (opcional, para base de datos y despliegue rápido)

### Paso 1: Configurar Variables de Entorno
Copia el archivo `.env.example` a `.env` y define tus credenciales:
```bash
cp .env.example .env
```

Ejemplo de configuración local:
```env
DB_USER=root
DB_PASSWORD=supersecreto
DB_HOST=localhost
DB_PORT=3306
DB_NAME=green_house_db
JWT_SECRET=tu_clave_secreta_aqui
PORT=8080
```

### Paso 2: Iniciar la Base de Datos (Docker)
Puedes levantar la base de datos MariaDB local rápidamente usando Docker Compose:
```bash
docker compose up -d db
```
El esquema de base de datos inicial se encuentra en `schema/green_house_db.sql`.

---

## 🧪 Pruebas Unitarias y Cobertura (Coverage)

Hemos implementado un conjunto robusto de pruebas unitarias cubriendo todas las capas del proyecto (Servicios, Handlers, Middleware, Seguridad y Repositorios). Para la capa de persistencia (repositorios), se utiliza `go-sqlmock` para simular la base de datos de manera aislada, rápida y determinista.

Para ejecutar los tests y validar la cobertura de código, utiliza los siguientes comandos:

### 1. Ejecutar todos los tests unitarios
Para ejecutar todas las pruebas en modo detallado (verbose):
```bash
go test -v ./...
```

### 2. Generar archivo de cobertura de código
Para correr las pruebas y generar un reporte de cobertura llamado `coverage.out`:
```bash
go test -coverprofile=coverage.out -coverpkg=./... ./...
```

### 3. Ver reporte de cobertura por consola
Para listar detalladamente el porcentaje de cobertura de cada función y archivo en tu terminal:
```bash
go tool cover -func=coverage.out
```

### 4. Ver reporte interactivo en el Navegador Web (HTML)
Para generar y abrir una página web interactiva local que resalta visualmente en verde las líneas de código cubiertas por los tests y en rojo las pendientes:
```bash
go tool cover -html=coverage.out
```

> [!NOTE]
> La meta del proyecto es mantener la cobertura agregada por encima del **80%**. Actualmente, la cobertura total del proyecto es de **91.0%**, con los servicios, handlers, middlewares y repositorios completamente probados.