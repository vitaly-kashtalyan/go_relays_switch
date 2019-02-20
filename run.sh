#!/usr/bin/env bash

docker build . -t go-relays-switch
docker run -i -t -p 8082:8082 --restart always go-relays-switch