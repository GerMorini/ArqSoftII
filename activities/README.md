# Activities-api

## Endpoints

verificar estado de la API

```bash
curl -i 'localhost:8081/healthz'
```
### Todos los endpoints estan en postman para testear.

listar actividades

```bash
curl -i 'localhost:8081/activities'
```

obtener actividad por su ID

```bash
curl -i 'localhost:8081/activities/64f1a6a1e4b0f1234567890a'
```

crear actividad (requiere JWT en Authorization)

```bash
TOKEN='...'
curl -i 'localhost:8081/activities' -X POST \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nombre": "Yoga Principiantes",
    "descripcion": "Clase suave para iniciar",
    "profesor": "12",               
    "diaSemana": "Lunes",
    "horaInicio": "09:00",
    "horaFin": "10:00",
    "capacidadMax": "25"
}'
```

actualizar actividad (requiere JWT en Authorization)

```bash
TOKEN='...'
ID='64f1a6a1e4b0f1234567890a'
curl -i "localhost:8081/activities/$ID" -X PUT \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "descripcion": "Clase actualizada",
    "horaFin": "10:30"
}'
```

eliminar actividad (requiere JWT en Authorization)

```bash
TOKEN='...'
ID='64f1a6a1e4b0f1234567890a'
curl -i "localhost:8081/activities/$ID" -X DELETE \
  -H "Authorization: Bearer $TOKEN"
```

consultar actividades a las que se está inscripto (requiere JWT de usuario no admin)

```bash
TOKEN='...'
ID=2
curl -i "localhost:8081/inscriptions/$ID" -X POST \
  -H "Authorization: Bearer $TOKEN"
```

inscribirse a una actividad (requiere JWT de usuario no admin)

```bash
TOKEN='...'
ID='64f1a6a1e4b0f1234567890a'
curl -i "localhost:8081/activities/$ID/inscribir" -X POST \
  -H "Authorization: Bearer $TOKEN"
```

desinscribirse de una actividad (requiere JWT de usuario no admin)

```bash
TOKEN='...'
ID='64f1a6a1e4b0f1234567890a'
curl -i "localhost:8081/activities/$ID/desinscribir" -X POST \
  -H "Authorization: Bearer $TOKEN"
```

> Nota: el puerto por defecto es 8080. Se puede cambiar con la variable `PORT_ACTIVIDADES_API`.

## Rápido (Docker Compose)

Si usas el repo con Docker Compose (recomendado para pruebas locales):

1. Desde la raíz del repo ejecutar:

```powershell
docker compose -f .\docker-compose.yml up -d --build
```

2. Espera a que el contenedor de Mongo esté listo antes de que `activities-api` intente conectarse. Revisa logs:

```powershell
docker compose -f .\docker-compose.yml logs -f mongo_activities_api
docker compose -f .\docker-compose.yml logs -f activities-api
```

En la configuración de compose `activities` expone por defecto el puerto `8081` en el host (para evitar conflicto con `users-api` que usa `8080`). Ajusta `baseUrl` según tu compose si es necesario.

## Autenticación (resumen)

Los endpoints protegidos requieren la cabecera HTTP:

  Authorization: Bearer <TOKEN>

El token lo obtienes con `POST /login` del servicio `users` y tiene 30 minutos de validez.

Claims incluidos en el token (generado por `users`):

- iss: "users-api"
- exp: expiración (Unix timestamp)
- id_usuario
- nombre
- apellido
- username
- email
- is_admin

Reglas específicas en Activities:

- `POST /activities`, `PUT /activities/:id`, `DELETE /activities/:id` requieren token válido.
- `POST /activities/:id/inscribir` y `POST /activities/:id/desinscribir` requieren token válido y que `is_admin` sea `false`.

## Postman / pruebas

He incluido una colección Postman lista para importar en:

`postman_collections/activities_postman_collection.json`

Importa la colección y ejecuta: primero `Users / Register` (opcional), luego `Users / Login` (guarda el token), y después las peticiones de `Activities` (Create guarda `activityId`).

## Notas y recomendaciones

- No guardes secretos (p.ej. `JWT_SECRET`) en repositorios públicos. Para entornos de producción usa Docker secrets o un gestor de secretos.
- En los contenedores, `MONGO_URI` debe apuntar al servicio de Compose (p.ej. `mongodb://mongo_activities_api:27017`), NO a `localhost`.
- Si cambias las imágenes base en los Dockerfiles a `latest`, ten en cuenta que esto puede alterar reproducibilidad; en producción prefiero tags fijos.

Si quieres que añada ejemplos automáticos (scripts bash/PowerShell) para ejecutar la secuencia completa (register → login → create → inscribir → desinscribir → delete), lo agrego en `tools/`.

## Autenticación

Los endpoints protegidos requieren cabecera HTTP `Authorization: Bearer <TOKEN>`.
El token JWT debe estar firmado con el secreto configurado en `JWT_SECRET`.

## Claims del token JWT

Datos generales:

- **iss**: emisor del token (users-api)
- **exp**: tiempo de expiración de 30 min

Datos del usuario:

- **id_usuario**
- **nombre**
- **apellido**
- **username**
- **email**
- **is_admin**

Reglas específicas en Activities:

- Para `POST /activities`, `PUT /activities/:id`, `DELETE /activities/:id`: requiere token válido.
- Para `POST /activities/:id/inscribir` y `POST /activities/:id/desinscribir`: requiere token válido y que `is_admin` sea `false`.

## Variables de entorno

- `PORT_ACTIVIDADES_API`: puerto de la API (por defecto `8080`).
- `MONGO_URI`: URL de conexión a MongoDB (por defecto `mongodb://localhost:27017`).
- `MONGO_DB`: nombre de la base de datos (por defecto `demo`).
- `JWT_SECRET`: secreto HMAC para validar tokens JWT (obligatorio).

## Comandos útiles

Generar un JWTSecret de 512 bits para firmar tokens:

```bash
openssl rand -hex 64
```

Decodificar partes del token JWT:

```bash
echo -n 'TOKEN' | base64 -d
```

Conectarse a MongoDB (adaptar al entorno):

```bash
# Si usa contenedor
docker exec -ti <mongo-container> mongosh

# O localmente
mongosh "mongodb://localhost:27017"
```


