#!/usr/bin/env bash

PROJECT="zlr-ca"
TAG="${PROJECT}:current"

DB="${PROJECT}-psql"
DB_TAG="${DB}:current"

echo "[DOCKER] stopping and removing old containers..."
docker rm $(docker stop $(docker ps -a -q --filter name=${PROJECT}$ --format="{{.ID}}"))
docker rm $(docker stop $(docker ps -a -q --filter name=${DB}$ --format="{{.ID}}"))
echo "done"
echo ""

if [ "$1" == "force" ]
then
    echo "[DOCKER] removing old image..."
    docker rmi ${TAG}
    docker rmi ${DB_TAG}
    echo "done"
    echo ""
fi

echo "[DOCKER-COMPOSE] up..."
docker-compose up -d
echo "done"
echo ""