#!/usr/bin/env bash

PROJECT="zlr-ca"
TAG="${PROJECT}:current"

echo "[DOCKER] stopping and removing old container..."
docker rm $(docker stop $(docker ps -a -q --filter ancestor=${TAG} --format="{{.ID}}"))
echo "done"
echo ""

if [ "$1" == "force" ]
then
    echo "[DOCKER] removing old image..."
    docker rmi ${TAG}
    echo "done"
    echo ""
fi

echo "[DOCKER] running docker build..."
docker build --rm=true -t ${TAG} .
echo "done"
echo ""

#echo "[DOCKER] cleanup dangling/intermediate images..."
#docker rmi $(docker images -f "dangling=true" -q)
#echo "done"
#echo ""

#docker run -d --name ${PROJECT} -p 8080:8080 --net=deployments_default --link zlr-ca-psql:current ${TAG} --network=deployments_default
