### BaseApp deployment script to create a new LXC Container via Docker
###
### Docker: http://www.docker.io

FROM ubuntu:latest
MAINTAINER Rich Tibbett <rich.tibbett@gmail.com>

# Update the base system
RUN apt-get update

# Tell debconf to run in non-interactive mode
ENV DEBIAN_FRONTEND noninteractive

# Install System Dependencies
RUN apt-get -y install build-essential golang git-core mercurial mysql-client mysql-server nginx pwgen python-setuptools vim-tiny

# Setup Go
RUN mkdir /go
ENV GOPATH  /go
ENV PATH $PATH:$GOPATH/bin

# Install Supervisord
RUN /usr/bin/easy_install supervisor
RUN /usr/bin/easy_install supervisor-stdout

# Get BaseApp Dependencies
#
# Annoyingly, we can not use `go get baseapp/...` because
# references to revel/app/routes package fail
RUN go get -v github.com/revel/revel github.com/revel/cmd/revel github.com/robfig/cron github.com/coopernurse/gorp code.google.com/p/go.crypto/bcrypt github.com/mattn/go-sqlite3 github.com/go-sql-driver/mysql github.com/ftrvxmtrx/gravatar github.com/russross/blackfriday

# Add Nginx frontend host
ADD ./docker/nginx_baseapp.vhost /etc/nginx/sites-available/default

# Stage BaseApp
ENV BASEAPP_PATH github.com/richtr/baseapp
ADD . $GOPATH/src/$BASEAPP_PATH

# Setup Nginx
RUN sed -i -e"s/keepalive_timeout\s*65/keepalive_timeout 2/" /etc/nginx/nginx.conf
RUN sed -i -e"s/keepalive_timeout 2/keepalive_timeout 2;\n\tclient_max_body_size 100m/" /etc/nginx/nginx.conf
RUN echo "daemon off;" >> /etc/nginx/nginx.conf

# Setup Supervisord
RUN cp $GOPATH/src/$BASEAPP_PATH/docker/supervisord.conf /etc/supervisord.conf

# Set start script permissions
RUN cp $GOPATH/src/$BASEAPP_PATH/docker/start.sh /start.sh
RUN chmod 755 /start.sh

# Expose Web Frontend (nginx) port only
EXPOSE 80

# Start required services when docker is instantiated
ENTRYPOINT ["/bin/bash", "/start.sh"]
