#!/bin/bash
# -----------------------------------------------------------------
# Script para EXPORTAR bases de datos (MySQL, MongoDB, Solr)
# -----------------------------------------------------------------
set -e # Salir inmediatamente si un comando falla

# --- Configuración ---
# Ajusta estos valores según tu entorno
MYSQL_CONTAINER="mysql-users-api"
MYSQL_USER="root"
MYSQL_PASS="root"
MYSQL_DB="users"

MONGO_CONTAINER="mongo-activities-api"
MONGO_DB="activities"

# --- Fin Configuración ---

BACKUP_DIR="backups_$(date +%Y-%m-%d_%H-%M-%S)"
mkdir -p "$BACKUP_DIR"

echo "--- Iniciando exportación ---"
echo "Directorio de backup: $BACKUP_DIR"
echo ""

echo "Exportando MySQL (DB: $MYSQL_DB)..."
docker exec "$MYSQL_CONTAINER" mysqldump -u"$MYSQL_USER" -p"$MYSQL_PASS" "$MYSQL_DB" > "$BACKUP_DIR/${MYSQL_DB}_backup.sql"
echo "OK: MySQL exportado."
echo ""

echo "Exportando MongoDB (DB: $MONGO_DB)..."
docker exec "$MONGO_CONTAINER" mongodump --db="$MONGO_DB" --archive=/tmp/mongo_backup.gz --gzip
docker cp "$MONGO_CONTAINER:/tmp/mongo_backup.gz" "$BACKUP_DIR/mongo_backup.gz"
docker exec "$MONGO_CONTAINER" rm /tmp/mongo_backup.gz
echo "OK: MongoDB exportado."
echo ""

echo "--- Exportación completada ---"
