#!/bin/bash

PACKAGE_NAME=github.com/richtr/baseapp

# Configure BaseApp execution paths
BASEAPP_PATH=/baseapp

# Configure BaseApp configuration path
BASEAPP_CONF_FILE=$BASEAPP_PATH/conf/app.conf

# Grab the best app.conf file
if [ ! -f $BASEAPP_CONF_FILE ]; then
	cp "${BASEAPP_CONF_FILE}.default" "$BASEAPP_CONF_FILE"
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

# Setup environment variable app.conf database overrides

[ -n "$BASEAPP_DATABASE_DRIVER" ] && sed -i .original \
    -e "s/db\.driver *=.*/db.driver = ${BASEAPP_DATABASE_DRIVER}/g" \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_DB_IMPORT" ] && sed -i .original \
    -e "s/db\.import *=.*/db.import = ${BASEAPP_DB_IMPORT}/g" \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_DB_SPEC" ] && sed -i .original \
    -e "s/db\.spec *=.*/db.spec = ${BASEAPP_DB_SPEC}/g" \
    $BASEAPP_CONF_FILE

# Remove sed back-up files
if [ -f $BASEAPP_CONF_FILE.original ]; then
	rm -f $BASEAPP_CONF_FILE.original
fi

# Set default BASEAPP_RUN_LEVEL if not set as an environment variable
[ -z "$BASEAPP_RUN_LEVEL" ] && BASEAPP_RUN_LEVEL=test

# Set default BASEAPP_RUN_PORT if not set as an environment variable
[ -z "$BASEAPP_RUN_PORT" ] && BASEAPP_RUN_PORT=9000

# Start Baseapp
echo "Running BaseApp at BASEAPP_RUN_LEVEL[${BASEAPP_RUN_LEVEL}]..."
mkdir -p $GOPATH/src/github.com/richtr && \
  cp -R $BASEAPP_PATH $GOPATH/src/$PACKAGE_NAME && \
  revel run $PACKAGE_NAME $BASEAPP_RUN_LEVEL $BASEAPP_RUN_PORT
