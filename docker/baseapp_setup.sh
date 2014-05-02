#!/bin/bash

# Configure BaseApp (MySQL Password, Database setup and Config setup)
if [ ! -f /baseapp-configured ]; then
	# General MySQL config changes
	sed -i -e "s/^bind-address\s*=\s*127.0.0.1/bind-address = 0.0.0.0/" /etc/mysql/my.cnf

	/usr/bin/mysqld_safe &
	sleep 10s

	# Generate new MySQL root password
	MYSQL_PASSWORD=`pwgen -c -n -1 12`
	echo mysql root password: $MYSQL_PASSWORD

	# Set MySQL root password
	mysqladmin -u root password $MYSQL_PASSWORD

	# Create required BaseApp databases
	mysqladmin -u root -p$MYSQL_PASSWORD CREATE baseapp
	mysqladmin -u root -p$MYSQL_PASSWORD CREATE baseapp_dev

	killall mysqld

	# Grab the best app.conf file
	BASEAPP_CONF_FILE=$GOPATH/src/$BASEAPP_PATH/conf/app.conf
	if [ ! -f $GOPATH/src/$BASEAPP_PATH/conf/app.conf ]; then
		BASEAPP_CONF_FILE=$GOPATH/src/$BASEAPP_PATH/conf/app.conf.default
	fi

	# Ensure app secret is set
	sed -i "s/<app_secret_please_change_me>/`pwgen -c -n -1 65`/g" $BASEAPP_CONF_FILE

	# Add MySQL password to BaseApp app.conf file
	sed -i "s/user:pass@tcp(localhost:3306)\/baseapp_dev?charset=utf8/root:$MYSQL_PASSWORD@tcp(localhost:3306)\/baseapp_dev?charset=utf8/g" $BASEAPP_CONF_FILE
	sed -i "s/user:pass@tcp(localhost:3306)\/baseapp?charset=utf8/root:$MYSQL_PASSWORD@tcp(localhost:3306)\/baseapp?charset=utf8/g" $BASEAPP_CONF_FILE

	if [ $BASEAPP_CONF_FILE == $GOPATH/src/$BASEAPP_PATH/conf/app.conf.default ]; then
		mv $BASEAPP_CONF_FILE $GOPATH/src/$BASEAPP_PATH/conf/app.conf
	fi

	touch /baseapp-configured
	sleep 10s

fi
