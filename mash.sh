#!/bin/bash

set -e

# --- Configurations --- 
export DEV_DATABASE_URL=postgres://pgadmin:pgpass@127.0.0.1:5432/game?sslmode=disable
export RELEASE_DATABASE_URL=

export RELEASE_MODE=dev 
export DEV_BOT_TOKEN="8122683105:AAHIVYGHFjgb0DCv1SHQYjv30Wg3olsnnxI" 

BASE_LOCATION=$(pwd)

generate_sqlc() {
	echo "Generating..."
	cd ./backend/
	sqlc generate
	cd ..
	echo "OK"	
}

build_frontend() {
	echo "Building frontend"
	cd ./frontend/
	npm run build
	cd ..
}

run_frontend() {
	echo "Running frontend dev server"
	cd ./frontend/
	npm run dev
	cd ..
}

test_build_frontend() {
	echo "Running frontend dev build server"
	cd ./frontend/
	npm run build
	serve -s dist -l 3000 
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
		cd ./backend/
		go run ./cmd/bot/.
		cd ..
		echo "The End"
		;;

	gensqlc)
		generate_sqlc
	;;
	fdev)
		run_frontend
	;;
	fbuild)
		build_frontend
	;;
	fserve)
		test_build_frontend
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
