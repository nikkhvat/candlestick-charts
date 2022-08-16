# Makefile

env-prepare: # create .env-file for secrets
	cp -n .example.env  .env

build: # build server
	go build -o ./.bin/app ./cmd/api/main.go

start: # start server
	./.bin/app

dev: # build and start server
	go build -o ./.bin/app ./cmd/api/main.go
	./.bin/app
