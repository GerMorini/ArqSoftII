# Search-api

Microservicio de búsqueda de actividades con caché de dos capas (local + distribuida) y motor Apache Solr.

## Endpoints

verificar estado de la API

```bash
curl -i 'localhost:8082/healthz'
```

buscar actividades con filtros

```bash
# Búsqueda básica sin filtros
curl -i 'localhost:8082/activities'

# Búsqueda con filtros por título
curl -i 'localhost:8082/activities?titulo=yoga'

# Búsqueda con filtros combinados
curl -i 'localhost:8082/activities?titulo=yoga&dia=Lunes&descripcion=principiantes'

# Búsqueda con paginación
curl -i 'localhost:8082/activities?page=1&count=10'

# Búsqueda completa
curl -i 'localhost:8082/activities?titulo=yoga&dia=Lunes&page=0&count=20'
```

> Nota: el puerto por defecto es 8080. Se puede cambiar con la variable `PORT_SEARCH_API`.

## Arquitectura

### Flujo de búsqueda

```
Usuario → API → Cache local (ccache) → Memcached → Solr (source of truth)
                      ↓                     ↓           ↓
                   60s TTL              60s TTL     RabbitMQ Consumer
                                                        ↓
                                          Create/Update/Delete → FlushAll()
```

### Componentes

1. **Cache local (ccache)**: Caché en memoria del proceso (TTL: 60s)
2. **Memcached**: Caché distribuida compartida entre instancias (TTL: 60s)
3. **Apache Solr**: Motor de búsqueda y fuente de verdad
4. **RabbitMQ Consumer**: Escucha eventos de actividades y mantiene Solr sincronizado

### Estrategia de invalidación

Cuando ocurre cualquier cambio en actividades (create/update/delete):

1. El evento se procesa desde RabbitMQ
2. Se actualiza/indexa/elimina en Solr
3. Se invalida **toda la caché** (FlushAll en ambas capas)

Esta estrategia simple garantiza consistencia eventual sin gestión compleja de claves.

## Rápido (Docker Compose)

Si usas el repo con Docker Compose (recomendado para pruebas locales):

1. Desde la raíz del repo ejecutar:

```bash
docker compose up -d --build
```

2. Espera a que Solr y Memcached estén listos. Revisa logs:

```bash
docker compose logs -f search-api
docker compose logs -f solr
docker compose logs -f memcached
```

En la configuración de compose `search` expone por defecto el puerto `8082` en el host.

## Variables de entorno

- `PORT_SEARCH_API`: puerto de la API (por defecto `8080`).
- `SOLR_HOST`: host del servidor Solr (por defecto `localhost`).
- `SOLR_PORT`: puerto del servidor Solr (por defecto `8983`).
- `SOLR_CORE`: nombre del core de Solr (por defecto `activities`).
- `MEMCACHED_HOST`: host del servidor Memcached (por defecto `localhost`).
- `MEMCACHED_PORT`: puerto del servidor Memcached (por defecto `11211`).
- `MEMCACHED_TTL_SECONDS`: TTL de la caché distribuida (por defecto `60`).
- `RABBITMQ_HOST`: host del servidor RabbitMQ (por defecto `localhost`).
- `RABBITMQ_PORT`: puerto del servidor RabbitMQ (por defecto `5672`).
- `RABBITMQ_USERNAME`: usuario de RabbitMQ (por defecto `guest`).
- `RABBITMQ_PASSWORD`: contraseña de RabbitMQ (por defecto `guest`).
- `RABBITMQ_QUEUE_NAME`: nombre de la cola de eventos (por defecto `activities_queue`).

## Comandos útiles

Ver documentos indexados en Solr:

```bash
curl 'http://localhost:8983/solr/demo/select?q=*:*&rows=100'
```

Buscar en Solr con filtros:

```bash
# Buscar por título
curl 'http://localhost:8983/solr/demo/select?q=titulo:yoga'

# Buscar por día de la semana
curl 'http://localhost:8983/solr/demo/select?q=dia_semana:Lunes'
```

Limpiar todos los documentos de Solr:

```bash
curl 'http://localhost:8983/solr/demo/update?commit=true' \
  -H 'Content-Type: application/json' \
  -d '{"delete":{"query":"*:*"}}'
```

Ver estadísticas de Memcached:

```bash
echo "stats" | nc localhost 11211
```

Limpiar toda la caché de Memcached:

```bash
echo "flush_all" | nc localhost 11211
```

Conectarse a MongoDB (desde contenedor):

```bash
docker exec -ti <mongo-container> mongosh
```

## Troubleshooting

### Problema: Resultados de búsqueda desactualizados

**Causa**: Caché aún no expiró o RabbitMQ consumer no procesó el evento.

**Solución**:
1. Verificar logs del consumer: `docker compose logs -f search-api`
2. Verificar que el evento llegó a RabbitMQ
3. Invalidar manualmente: `echo "flush_all" | nc localhost 11211`

### Problema: No se encuentran actividades recién creadas

**Causa**: Evento de RabbitMQ no procesado o error al indexar en Solr.

**Solución**:
1. Verificar logs del consumer
2. Verificar que Solr está corriendo: `curl http://localhost:8983/solr/demo/admin/ping`
3. Revisar documentos en Solr manualmente

### Problema: Error "connection refused" a Memcached/Solr

**Causa**: Servicios no están corriendo o variables de entorno mal configuradas.

**Solución**:
1. En Docker Compose, usar nombres de servicio (ej: `solr`, `memcached`) NO `localhost`
2. Verificar que los servicios están corriendo: `docker compose ps`
3. Revisar variables de entorno en `.env`

## RabbitMQ Consumer

El consumer procesa tres tipos de eventos:

### Evento CREATE

```json
{
  "action": "create",
  "id": "64f1a6a1e4b0f1234567890a",
  "nombre": "Yoga Principiantes",
  "descripcion": "Clase suave",
  "dia": "Lunes"
}
```

Indexa la actividad en Solr y limpia ambas cachés.

### Evento UPDATE

```json
{
  "action": "update",
  "id": "64f1a6a1e4b0f1234567890a",
  "nombre": "Yoga Avanzado",
  "descripcion": "Clase modificada",
  "dia": "Martes"
}
```

Reindexiza la actividad en Solr y limpia ambas cachés.

### Evento DELETE

```json
{
  "action": "delete",
  "id": "64f1a6a1e4b0f1234567890a",
  "nombre": "",
  "descripcion": "",
  "dia": ""
}
```

Elimina la actividad de Solr y limpia ambas cachés.

## Notas y recomendaciones

- El TTL de 60 segundos balancea consistencia y performance
- FlushAll garantiza consistencia eventual sin gestión compleja de claves
- Solr es la única fuente de verdad; las cachés son solo para optimización
- En producción considera monitored de Solr y Memcached
- Si necesitas búsquedas más complejas, aprovecha las capacidades de Solr (facets, highlighting, etc.)
