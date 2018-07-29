#!/usr/bin/env bash

#
# sometimes wait-for-it reports to fast that postgres is available..
#
./wait-for-it.sh -t 0 postgres:5432

#
# .. so waiting after that some seconds more..
#
sleep 5

#
# .. and then finally start the server
#

server -h ${DB_HOST} -pt ${DB_PORT} -u ${DB_USER} -p ${DB_PSWD} -d ${DB_NAME} -s ${DB_SCHEMA}