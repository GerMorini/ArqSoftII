# Users-api

## Endpoints

verificar estado de la API

```bash
curl -i 'localhost:8080/healthz'
```

registrar usuario

```bash
curl -i 'localhost:8080/register' -X POST -d '{
    "nombre": "Pepe",
    "apellido": "Gomez",
    "username": "pgomez31",
    "email": "pepe.gom@yahoo.com",
    "password": "secreto"
}'
```

loggearse por nombre de usuario (genera un token)

```bash
curl -i 'localhost:8080/login' -X POST -d '{
    "username": "pgomez31",
    "password": "secreto"
}'
```

loggearse por email (genera un token)

```bash
curl -i 'localhost:8080/login' -X POST -d '{
    "email": "pepe.gom@yahoo.com",
    "password": "secreto"
}'
```

obtener usuario por su ID

```bash
curl -i 'localhost:8080/users/1'
```

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

## Comandos útiles

Generar un JWTSecret de 512 bits para firmar tokens:

```bash
openssl rand -hex 64
```

Manipular la base de datos desde el contenedor:

```bash
docker exec -ti mysql-users-api mysql -u root -p users
```

Decodificar partes del token JWT:

```bash
echo -n 'TOKEN' | base64 -d
```
