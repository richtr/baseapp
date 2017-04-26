###
### ~*~ BaseApp ~*~
###
### A Bootstrap-based web application built in Go on top of the Revel Web Framework
###
### https://github.com/richtr/baseapp
###

FROM golang:alpine

MAINTAINER Rich Tibbett

# Set default BaseApp run level
ENV BASEAPP_RUN_LEVEL test

ENV BASEAPP_PATH $GOPATH/src/github.com/richtr/baseapp

# Stage BaseApp
ADD . $BASEAPP_PATH

# Add start script
ADD ./start.sh /start.sh

# Install baseapp golang dependencies and set start script permissions
RUN apk add --no-cache gcc g++ git pwgen bash && \
		cd $BASEAPP_PATH && \
		go get -v ./... github.com/revel/revel github.com/revel/cmd/revel && \
		chmod 755 /start.sh

# Expose BaseApp port
EXPOSE 9000

# Configure and start BaseApp on load
ENTRYPOINT ["/bin/bash", "/start.sh"]
