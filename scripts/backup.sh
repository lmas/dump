#!/bin/bash
set -e # exit script on errors

BACKUP_DIR="./dir"

PREFIX="test"
TIMESTAMP=$(date +"%Y%m%d-%H%M")

#############################################################################
# Create new backup
# TODO: add encryption

FILE="$BACKUP_DIR/$PREFIX-$TIMESTAMP.tar.bz2"

# check if dir exists
if [ ! -d "$BACKUP_DIR" ]; then
	mkdir "$BACKUP_DIR"
	echo "Created backup dir: $BACKUP_DIR"
fi

tar -cjf "$FILE" example/
echo "Create backup: $FILE"

#############################################################################
# Rotate old backups

find "$BACKUP_DIR/" -mmin +7 -exec echo {} \; | sort
find "$BACKUP_DIR/" -mmin +7 -exec rm {} \;



# backup every hour the first day
#find "$BACKUP_HOURLY/" -mtime +1 -exec rm {} \;

# backup all days the first week
#find "$BACKUP_DAILY/" -mtime +7 -exec rm {} \;

# backup all weeks
