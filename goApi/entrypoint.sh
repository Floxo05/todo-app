#!/bin/sh

FILE_PATH="/is_migrated"

until timeout 1 bash -c 'echo > /dev/tcp/db/3306'; do
  echo "Waiting for the Database to start..."
  sleep 1
done

if [ ! -f "$FILE_PATH" ]; then
    go run cmd/migrate/up/up.go

    touch $FILE_PATH
fi

exec "./main"