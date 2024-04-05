#!/usr/bin/env sh
set -eax

APP_NAME="scgb"

mkdir -p "cmd/$APP_NAME/generated"
echo "$APP_NAME" | tr -d '\n' > "cmd/$APP_NAME/generated/app-name"
pwd | tr -d '\n' > "cmd/$APP_NAME/generated/source-dir"

cd "./cmd/$APP_NAME"
go install
