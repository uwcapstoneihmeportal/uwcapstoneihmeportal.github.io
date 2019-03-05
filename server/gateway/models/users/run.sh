#!/usr/bin/env bash
#export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
export MYSQL_ROOT_PASSWORD="tester"
docker stop myuserstore
docker system prune -f
docker run -d -p 3306:3306 --name myuserstore \
 -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
 -e MYSQL_DATABASE=userstore \
 taehyun123/myuserstore