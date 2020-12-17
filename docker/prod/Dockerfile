####################
### Build JS app ###
####################
FROM node:alpine as web

RUN apk update && \
  apk upgrade --no-cache && \
  apk add --no-cache build-base python

RUN mkdir /app
COPY web /app
WORKDIR /app

RUN yarn
RUN yarn build


####################
### Build Go app ###
####################
FROM golang:alpine as app

RUN mkdir /app
COPY go.mod go.sum /app/
WORKDIR /app

RUN go mod download

COPY types /app/types
COPY internal /app/internal
COPY *.go /app/
COPY user_agents.txt /app/
RUN go build -o stock_scraper .


#####################
### Release image ###
#####################
FROM alpine
RUN mkdir /app
COPY --from=app /app/stock_scraper /app/stock_scraper
COPY --from=web /app/build /app/static
COPY user_agents.txt /app/
WORKDIR /app
