#!/usr/bin/env bash

BACKUP_DIR="/var/data/backups"
REAL_FILE="/home/leo/bdkursach/backend/scripts/backups"
DATE=$(date +%Y-%m-%d_%H-%M-%S)
CONTAINER="postgres_kurs"
USER="admin"
MAXSIZE=10
DB="db"

docker exec $CONTAINER pg_dump -U $USER -d $DB -F c -f $BACKUP_DIR/$DATE.dump
SIZE=$(ls $REAL_FILE | wc -l)
echo $SIZE
DIFF=$((MAXSIZE-SIZE))
echo $DIFF
if [ $DIFF -lt 1 ]; then
	DIFF=$((SIZE-MAXSIZE))
	ls -t $REAL_FILE | tail -n+$DIFF | xargs -I{} rm $REAL_FILE/{}
fi
