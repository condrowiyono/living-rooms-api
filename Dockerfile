FROM golang:1.12.0-alpine3.9

RUN apk add --no-cache bash ca-certificates git gcc libc-dev curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh


# install service dependencies
RUN mkdir -p /go/src/living_rooms
WORKDIR /go/src/github.com/condrowiyono/living-rooms-api

COPY ./Gopkg.toml .
COPY ./Gopkg.lock .
RUN dep ensure -v --vendor-only

COPY . .

RUN go build -o main .

CMD ["/go/src/github.com/condrowiyono/living-rooms-api/main"]