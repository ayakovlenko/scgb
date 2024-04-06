#!/usr/bin/env sh
set -eax

APP_NAME="scgb"

mkdir -p "cmd/$APP_NAME/generated"
echo "$APP_NAME" | tr -d '\n' > "cmd/$APP_NAME/generated/app-name"
pwd | tr -d '\n' > "cmd/$APP_NAME/generated/source-dir"
touch "cmd/$APP_NAME/generated/app-hash"

cd "./cmd/$APP_NAME"
go install
