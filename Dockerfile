###
### ~*~ BaseApp ~*~
###
### A Bootstrap-based web application built in Go on top of the Revel Web Framework
###
### https://github.com/richtr/baseapp
###

FROM golang:1.8-alpine

MAINTAINER Rich Tibbett

# Install baseapp dependencies
RUN apk add --no-cache gcc g++ git bash perl

RUN go get github.com/revel/revel && \
    go get github.com/revel/cmd/revel

# Set default BaseApp environment variables
ENV BASEAPP_RUN_LEVEL test
ENV BASEAPP_SERVER_PORT 9000
ENV BASEAPP_PATH github.com/richtr/baseapp
ENV BASEAPP_DIR $GOPATH/src/$BASEAPP_PATH

# Add BaseApp
ADD . $BASEAPP_DIR

# Expose BaseApp port
EXPOSE 9000

# Configure and start BaseApp on load
WORKDIR $BASEAPP_DIR
ENTRYPOINT ["/bin/bash", "start.sh"]
