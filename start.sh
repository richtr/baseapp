#!/bin/bash

# Configure BaseApp execution paths
BASEAPP_DIR=github.com/richtr/baseapp
BASEAPP_PATH=$GOPATH/src/$BASEAPP_DIR

# Configure BaseApp configuration path
BASEAPP_CONF_FILE=$BASEAPP_PATH/conf/app.conf

# Grab the best app.conf file
if [ ! -f $BASEAPP_CONF_FILE ]; then
	cp $BASEAPP_CONF_FILE.default $BASEAPP_CONF_FILE
fi

# Generate app secret (better than pwgen as this works cross-platform)
APP_SECRET=`env LC_CTYPE=C tr -dc "a-zA-Z0-9-_\$\?" < /dev/urandom | head -c 65`

# Ensure app secret is set
sed -i .original \
    -e "s/app\.secret *= *<app_secret_please_change_me>/app.secret = ${APP_SECRET}/g" \
    $BASEAPP_CONF_FILE

# Expose on all available network interfaces exposed to Docker container
sed -i .original \
    -e "s/http\.addr *= *localhost/http.addr =/g" \
    $BASEAPP_CONF_FILE

# Remove sed back-up files
if [ -f $BASEAPP_CONF_FILE.original ]; then
	rm -f $BASEAPP_CONF_FILE.original
fi

# Install Revel
echo "Installing Revel..."
go get -v github.com/revel/revel github.com/revel/cmd/revel

# Install and start Baseapp (in test mode)
echo "Installing and running BaseApp..."
cd $BASEAPP_PATH && \
   go get -v ./... && \
   cd $GOPATH/src && \
   revel run $BASEAPP_DIR test
