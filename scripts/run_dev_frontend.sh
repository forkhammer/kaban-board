#!/bin/bash

set -o allexport
source .env
set +o allexport

cd ./frontend/board
npm run start