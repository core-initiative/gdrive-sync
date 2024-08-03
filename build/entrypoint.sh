#!/bin/bash -e

APP_ENV=${APP_ENV:-local}

echo "[`date`] Running entrypoint script in the '${APP_ENV}' environment..."

CONFIG_FILE=./config.yaml

echo "[`date`] Starting server..."
./server -config ${CONFIG_FILE}