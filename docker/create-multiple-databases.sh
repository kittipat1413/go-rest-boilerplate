#!/bin/bash

set -e
set -u

function create_user_and_database() {
	local database=$1
	local owner=$2
	echo "Creating database '$database' and grant all privileges to '$owner'"
	psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
	    CREATE DATABASE "$database";
	    GRANT ALL PRIVILEGES ON DATABASE "$database" TO "$owner";
EOSQL
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ] && [ -n "$POSTGRES_USER" ]; then
	echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
		create_user_and_database $db $POSTGRES_USER
	done
	echo "Multiple databases created"
fi