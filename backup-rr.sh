#!/bin/bash

# Get current timestamp
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_ROOT="./backups"
BACKUP_DIR="$BACKUP_ROOT/$TIMESTAMP"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Load environment variables from .env
if [ -f .env ]; then
  # Auto-export variables
  set -a
  source .env
  set +a
else
  echo "Error: .env file not found."
  exit 1
fi

echo "Starting backup to $BACKUP_DIR..."

# --- 1. Backup Databases ---
# We use docker compose exec -T (disable TTY allocation) to pipe stdout
echo "dumping database: $PANSMAN_DB_NAME..."
docker compose exec -T db mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" "$PANSMAN_DB_NAME" > "$BACKUP_DIR/${PANSMAN_DB_NAME}.sql"

echo "dumping database: $MANOSMAN_DB_NAME..."
docker compose exec -T db mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" "$MANOSMAN_DB_NAME" > "$BACKUP_DIR/${MANOSMAN_DB_NAME}.sql"

echo "dumping database: $ANPANSMAN_DB_NAME..."
docker compose exec -T db mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" "$ANPANSMAN_DB_NAME" > "$BACKUP_DIR/${ANPANSMAN_DB_NAME}.sql"


# --- 2. Backup wp-content directories ---
echo "Archiving wp-content for anpan.sman.cloud..."
if [ -d "anpan_sman_cloud/wp-content" ]; then
    tar -czf "$BACKUP_DIR/anpan_sman_cloud.tar.gz" -C anpan_sman_cloud wp-content
else
    echo "Warning: anpan_sman_cloud/wp-content not found."
fi

echo "Archiving wp-content for pan.sman.cloud..."
if [ -d "pan_sman_cloud/wp-content" ]; then
    tar -czf "$BACKUP_DIR/pan_sman_cloud.tar.gz" -C pan_sman_cloud wp-content
else
    echo "Warning: pan_sman_cloud/wp-content not found."
fi

echo "Archiving wp-content for mano.sman.cloud..."
if [ -d "mano_sman_cloud/wp-content" ]; then
    tar -czf "$BACKUP_DIR/mano_sman_cloud.tar.gz" -C mano_sman_cloud wp-content
else
    echo "Warning: mano_sman_cloud/wp-content not found."
fi

# --- 3. Backup Configuration Files ---
echo "Backing up configuration files..."
cp docker-compose.yml "$BACKUP_DIR/"
cp .env "$BACKUP_DIR/"
cp README.md "$BACKUP_DIR/"
if [ -f "vozigo_wp.nginx" ]; then
    cp vozigo_wp.nginx "$BACKUP_DIR/"
fi

echo "------------------------------------------------"
echo "Backup finished successfully."
echo "Location: $BACKUP_DIR"
ls -lh "$BACKUP_DIR"

# --- 4. Upload to Google Drive ---

echo "Uploading backup to Google Drive..."
# rclone copy -P --drive-chunk-size 64M "$BACKUP_DIR" "gdrive:backup_vozigo_wp/$BACKUP_DIR"
rclone copy --drive-chunk-size 64M backups/ gdrive:backups/vozigo_wp/
