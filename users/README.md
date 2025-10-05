# Users-api

endpoints:

```bash
# verificar estado de la API
curl -i 'localhost:8080/healthz'

# registrar usuario
curl -i 'localhost:8080/register' -X POST -d '{
    "nombre": "Pepe",
    "apellido": "Gomez",
    "username": "pgomez31",
    "email": "pepe.gom@yahoo.com",
    "password": "secreto"
}'

# loggearse por nombre de usuario (genera un token)
curl -i 'localhost:8080/login' -X POST -d '{
    "username": "pgomez31",
    "password": "secreto"
}'

# loggearse por email (genera un token)
curl -i 'localhost:8080/login' -X POST -d '{
    "email": "pepe.gom@yahoo.com",
    "password": "secreto"
}'
```

generar un JWTSecret de 512 bits para firmar tokens:

```bash
openssl rand -hex 64

```

manipular la base de datos desde el contenedor:

```bash
docker exec -ti mysql-users-api mysql -u root -p users
```

decodificar partes del token JWT:
```bash
echo -n 'TOKEN' | base64 -d
```