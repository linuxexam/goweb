#!/bin/env bash
./build-app.sh
docker build -t goweb -f ./Dockerfile ../
