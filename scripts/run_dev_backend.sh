#!/bin/bash

set -o allexport
source .env
set +o allexport

cd ./backend
go run main