#!/bin/bash
# -----------------------------------------------------------------
# Script para IMPORTAR bases de datos (MySQL, MongoDB, Solr)
# -----------------------------------------------------------------
set -e # Salir inmediatamente si un comando falla

# --- Configuración ---
# Ajusta estos valores según tu entorno
MYSQL_CONTAINER="mysql-users-api"
MYSQL_USER="root"
MYSQL_PASS="root"
MYSQL_DB="users"

ACTIVITIES_CONTAINER="activities-api"
MONGO_CONTAINER="mongo-activities-api"
MONGO_DB="activities" # Necesario para --drop

# --- Fin Configuración ---

# Encontrar el directorio de backup más reciente
BACKUP_DIR=$(find . -maxdepth 1 -type d -name "backups_*" | sort -r | head -n 1)

if [ -z "$BACKUP_DIR" ]; then
    echo "Error: No se encontró ningún directorio 'backups_*'."
    echo "Asegúrate de ejecutar el script de exportación primero."
    exit 1
fi

echo "--- Iniciando importación desde $BACKUP_DIR ---"
echo ""

# 1. MySQL
MYSQL_FILE="$BACKUP_DIR/${MYSQL_DB}_backup.sql"
if [ -f "$MYSQL_FILE" ]; then
    echo "Importando MySQL (DB: $MYSQL_DB)..."
    docker exec -i "$MYSQL_CONTAINER" mysql -u"$MYSQL_USER" -p"$MYSQL_PASS" "$MYSQL_DB" < "$MYSQL_FILE"
    echo "OK: MySQL importado."
else
    echo "Advertencia: No se encontró archivo de backup para MySQL ($MYSQL_FILE)"
fi
echo ""

# 2. MongoDB
MONGO_FILE="$BACKUP_DIR/mongo_backup.gz"
if [ -f "$MONGO_FILE" ]; then
    echo "Importando MongoDB (DB: $MONGO_DB)..."
    docker cp "$MONGO_FILE" "$MONGO_CONTAINER:/tmp/mongo_backup.gz"
    # Usar --drop para borrar colecciones existentes antes de importar
    docker exec "$MONGO_CONTAINER" mongorestore --db="$MONGO_DB" --archive=/tmp/mongo_backup.gz --gzip --drop
    docker exec "$MONGO_CONTAINER" rm /tmp/mongo_backup.gz
    echo "OK: MongoDB importado."
else
    echo "Advertencia: No se encontró archivo de backup para MongoDB ($MONGO_FILE)"
fi

echo "Indexando solr..."
docker exec -ti $ACTIVITIES_CONTAINER reindex
echo "OK: Actividades reindexadas correctamente"

echo ""
echo "--- Importación completada ---"
