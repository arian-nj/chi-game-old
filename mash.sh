#!/bin/bash

set -e

# --- Configurations --- 
export DEV_DATABASE_URL=postgres://pgadmin:pgpass@127.0.0.1:5432/game?sslmode=disable
export RELEASE_DATABASE_URL=

export RELEASE_MODE=dev 
export DEV_BOT_TOKEN="8027085911:AAGVhwYk7erGCW4CYPBxeCXsPXZqX8w4SPM" 

BASE_LOCATION=$(pwd)

generate_sqlc() {
	echo "Generating..."
	sqlc generate
	echo "OK"	
}

build_frontend() {
	echo "Building frontend"
	cd ./frontend/
	npm run build
	cd ..
}

COMMAND=$1

if [ -z "COMMAND" ]; then
	echo "i want more commands"
fi

case $COMMAND in
	front)
		build_frontend
		;;
	run)
		generate_sqlc
		echo "Starting Core API"
		go run ./cmd/bot/.

		echo "The End"
		;;

	gensqlc)
		generate_sqlc
	;;
	dev)
		echo "serving dev frontend"
		cd ./frontend/
		npm run dev
	;;

	treafik)
		sudo traefik \
	--entrypoints.web.address=:8080 \
	--providers.file.filename=treafik.yaml \
	--log.level=DEBUG
	;;
	*)
		echo "command '$COMMAND' is Unknown"
esac

exit 0
