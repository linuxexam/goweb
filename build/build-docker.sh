#!/bin/env bash
# run this script under repo root

self_path=$(cd `dirname $0`; pwd)
docker build -t goweb -f $self_path/Dockerfile .
