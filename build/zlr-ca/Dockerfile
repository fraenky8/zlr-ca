FROM golang:1.10.3

ENV DB_HOST=192.168.99.100
ENV DB_PORT=5432ENV
ENV DB_USER=postgres
ENV DB_PSWD=mysecretpassword
ENV DB_NAME=postgres
ENV DB_SCHEMA=zlr_ca

ARG project=github.com/fraenky8/zlr-ca
ARG codedir=./src/$project

COPY wait-for-it.sh /wait-for-it.sh
COPY wait-for-it-wrapper.sh /wait-for-it-wrapper.sh
RUN chmod +x /*.sh

#
# setup project and dependencies
#
RUN go get -u -v $project/...

WORKDIR $codedir/cmd/server
RUN go install

WORKDIR /

#
# start the rest api server
#
CMD server -h ${DB_HOST} -pt ${DB_PORT} -u ${DB_USER} -p ${DB_PSWD} -d ${DB_NAME} -s ${DB_SCHEMA}

EXPOSE 8080