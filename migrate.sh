#!/bin/bash

DB_URL="postgres://postgres:postgres@localhost:5432/apiserver?sslmode=disable"

case $1 in
  create)
    if [ -z "$2" ]; then
      echo "Usage: ./migrate.sh create migration_name"
      exit 1
    fi
    docker run --rm -v $(pwd)/migrations:/migrations \
      migrate/migrate:v4.17.0 \
      create -ext sql -dir /migrations -seq $2
    echo "Created migration: $2"
    echo "Edit files in migrations/ folder"
    ;;
  
  up)
    docker run --rm -v $(pwd)/migrations:/migrations \
      --network host \
      migrate/migrate:v4.17.0 \
      -path=/migrations -database "$DB_URL" up
    echo "Migrations applied!"
    ;;
  
  down)
    docker run --rm -v $(pwd)/migrations:/migrations \
      --network host \
      migrate/migrate:v4.17.0 \
      -path=/migrations -database "$DB_URL" down 1
    echo "Last migration rolled back!"
    ;;
  
  version)
    docker run --rm --network host \
      migrate/migrate:v4.17.0 \
      -database "$DB_URL" version
    ;;
  
  *)
    echo "Usage: ./migrate.sh {create|up|down|version}"
    echo ""
    echo "Commands:"
    echo "  create <name>  - Create new migration"
    echo "  up            - Apply all pending migrations"
    echo "  down          - Rollback last migration"
    echo "  version       - Check current version"
    exit 1
    ;;
esac
