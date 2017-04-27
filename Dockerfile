###
### ~*~ BaseApp ~*~
###
### A Bootstrap-based web application built in Go on top of the Revel Web Framework
###
### https://github.com/richtr/baseapp
###

FROM golang:1.8-alpine

MAINTAINER Rich Tibbett

# Install baseapp golang dependencies and set start script permissions
RUN apk add --no-cache gcc g++ git bash

RUN go get github.com/revel/revel && \
    go get github.com/revel/cmd/revel

# Set default BaseApp environment variables
ENV BASEAPP_RUN_LEVEL test
ENV BASEAPP_RUN_PORT 9000
ENV BASEAPP_PATH /baseapp

# Add BaseApp
ADD . $BASEAPP_PATH

# Expose BaseApp port
EXPOSE 9000

# Configure and start BaseApp on load
ENTRYPOINT ["/bin/bash", "/baseapp/start.sh"]
