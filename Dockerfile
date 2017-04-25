###
### ~*~ BaseApp ~*~
###
### A Bootstrap-based web application built in Go on top of the Revel Web Framework
###
### https://github.com/richtr/baseapp
###

FROM golang:alpine

MAINTAINER Rich Tibbett

ENV BASEAPP_PATH $GOPATH/src/github.com/richtr/baseapp

WORKDIR $BASEAPP_PATH

# Stage BaseApp
ADD . $BASEAPP_PATH

# Add start script
ADD ./start.sh /start.sh

# Install baseapp golang dependencies and set start script permissions
RUN apk add --no-cache gcc g++ git pwgen bash && \
		go get github.com/revel/revel && \
		go get github.com/revel/cmd/revel && \
		go get ./... && \
		chmod 755 /start.sh

# Expose BaseApp port
EXPOSE 9000

# Configure and start BaseApp on load
ENTRYPOINT ["/bin/bash", "/start.sh"]
