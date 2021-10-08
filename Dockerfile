ARG BUILD_FROM="rodrigodiez/golang-cron"
FROM $BUILD_FROM

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ONBUILD COPY . /go/src/app
ONBUILD RUN go get -v -d
ONBUILD RUN go install -v
