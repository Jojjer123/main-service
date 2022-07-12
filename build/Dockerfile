FROM golang:1.18 as build

ENV MAIN_ROOT=/go/src/main-service
ENV CGO_ENABLED=0

RUN mkdir -p $MAIN_ROOT/

COPY . $MAIN_ROOT

RUN cd $MAIN_ROOT && GO111MODULE=on go build -o /go/bin/main ./


FROM alpine:3.11
RUN apk add bash
ENV HOME=/home/main-service
RUN mkdir $HOME
WORKDIR $HOME

COPY . $HOME

COPY --from=build /go/bin/main /usr/local/bin/

EXPOSE 8080

CMD ["main"]