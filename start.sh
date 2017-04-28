#!/bin/bash

# Configure BaseApp execution paths
[ -z "$BASEAPP_PATH" ] && BASEAPP_PATH=github.com/richtr/baseapp
[ -z "$BASEAPP_DIR" ] && BASEAPP_DIR=$GOPATH/src/$BASEAPP_PATH

# Set default BASEAPP_RUN_LEVEL if not set as an environment variable
[ -z "$BASEAPP_RUN_LEVEL" ] && BASEAPP_RUN_LEVEL=test

# Set default BASEAPP_SERVER_PORT if not set as an environment variable
[ -z "$BASEAPP_SERVER_PORT" ] && BASEAPP_SERVER_PORT=9000

# Configure BaseApp configuration path
BASEAPP_CONF_FILE=$BASEAPP_DIR/conf/app.conf

# Grab the best app.conf file
if [ ! -f $BASEAPP_CONF_FILE ]; then
  cp "${BASEAPP_CONF_FILE}.default" "$BASEAPP_CONF_FILE"
fi

# Generate app secret (better than pwgen as this works cross-platform)
APP_SECRET=`env LC_CTYPE=C tr -dc "a-zA-Z0-9-_\$\?" < /dev/urandom | head -c 65`

# Ensure app secret is set
perl -i -0777 -0pe \
    's#app\.secret( *)= *<app_secret_please_change_me>#app.secret$1= '"${APP_SECRET}"'#g' \
    $BASEAPP_CONF_FILE

#### SETUP BASEAPP ENVIRONMENT VARIABLE OVERRIDES ####

# APPLICATION OVERRIDES

[ -n "$BASEAPP_APPLICATION_NAME" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)app\.name( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1app.name$2= "'"${BASEAPP_APPLICATION_NAME}"'"$3#s' \
    $BASEAPP_CONF_FILE

# SERVER OVERRIDES

[ -n "$BASEAPP_SERVER_HOST" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)http\.host( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1http.host$2= "'"${BASEAPP_SERVER_HOST}"'"$3#s' \
    $BASEAPP_CONF_FILE

perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)http\.addr( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1http.addr$2= '"${BASEAPP_SERVER_ADDR}"'$3#s' \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_SERVER_PORT" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)http\.port( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1http.port$2= '"${BASEAPP_SERVER_PORT}"'$3#s' \
    $BASEAPP_CONF_FILE


if [ -n "$BASEAPP_SERVER_SSLCERT" && -n "$BASEAPP_SERVER_SSLKEY" ]; then

  perl -i -0777 -0pe \
      's#['"${BASEAPP_RUN_LEVEL}"'](.*)http\.ssl( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1http.ssl$2= true$3#s' \
      $BASEAPP_CONF_FILE

  perl -i -0777 -0pe \
      's#['"${BASEAPP_RUN_LEVEL}"'](.*)http\.sslcert( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1http.sslcert$2= "'"${BASEAPP_SERVER_SSLCERT}"'"$3#s' \
      $BASEAPP_CONF_FILE

  perl -i -0777 -0pe \
      's#['"${BASEAPP_RUN_LEVEL}"'](.*)http\.sslkey( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1http.sslkey$2= "'"${BASEAPP_SERVER_SSLKEY}"'"$3#s' \
      $BASEAPP_CONF_FILE

fi

# DATABASE OVERRIDES

[ -n "$BASEAPP_DATABASE_DRIVER" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)db\.driver( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1db.driver$2= "'"${BASEAPP_DATABASE_DRIVER}"'"$3#s' \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_DATABASE_IMPORT" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)db\.import( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1db.import$2= "'"${BASEAPP_DATABASE_IMPORT}"'"$3#s' \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_DATABASE_SPEC" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)db\.spec( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1db.spec$2= "'"${BASEAPP_DATABASE_SPEC}"'"$3#s' \
    $BASEAPP_CONF_FILE

# EMAIL OVERRIDES

[ -n "$BASEAPP_MAILER_SERVER" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)mailer\.server( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1mailer.server$2= '"${BASEAPP_MAILER_SERVER}"'$3#s' \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_MAILER_PORT" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)mailer\.port( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1mailer.port$2= '"${BASEAPP_MAILER_PORT}"'$3#s' \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_MAILER_USERNAME" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)mailer\.username( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1mailer.username$2= '"${BASEAPP_MAILER_USERNAME}"'$3#s' \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_MAILER_PASSWORD" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)mailer\.password( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1mailer.password$2= "'"${BASEAPP_MAILER_PASSWORD}"'"$3#s' \
    $BASEAPP_CONF_FILE

[ -n "$BASEAPP_MAILER_FROMADDRESS" ] && perl -i -0777 -0pe \
    's#['"${BASEAPP_RUN_LEVEL}"'](.*)mailer\.fromaddress( *)=[^\r\n]*(.*)#['"${BASEAPP_RUN_LEVEL}"']$1mailer.fromaddress$2= '"${BASEAPP_MAILER_FROMADDRESS}"'$3#s' \
    $BASEAPP_CONF_FILE

#### START BASEAPP ####

echo "Starting BaseApp in ${BASEAPP_RUN_LEVEL} mode..."
revel run $BASEAPP_PATH $BASEAPP_RUN_LEVEL $BASEAPP_SERVER_PORT
