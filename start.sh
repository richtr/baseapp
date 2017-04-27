#!/bin/bash

# Configure BaseApp execution paths
[ -z "$BASEAPP_PATH" ] && BASEAPP_PATH=github.com/richtr/baseapp
[ -z "$BASEAPP_DIR" ] && BASEAPP_DIR=$GOPATH/src/$BASEAPP_PATH

# Configure BaseApp configuration path
BASEAPP_CONF_FILE=$BASEAPP_DIR/conf/app.conf

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
echo "Running BaseApp in ${BASEAPP_RUN_LEVEL} mode..."
revel run $BASEAPP_PATH $BASEAPP_RUN_LEVEL $BASEAPP_RUN_PORT
