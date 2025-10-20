# Activities-api

## Endpoints

verificar estado de la API

```bash
curl -i 'localhost:8081/healthz'
```

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
    "dia_semana": "Lunes",
    "hora_inicio": "09:00",
    "hora_fin": "10:00",
    "capacidad_max": "25"
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
    "hora_fin": "10:30"
}'
```

eliminar actividad (requiere JWT en Authorization)

```bash
TOKEN='...'
ID='64f1a6a1e4b0f1234567890a'
curl -i "localhost:8081/activities/$ID" -X DELETE \
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


