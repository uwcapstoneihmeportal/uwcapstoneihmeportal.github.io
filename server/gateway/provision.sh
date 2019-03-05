#!/usr/bin/env bash
export TLSCERT=/etc/letsencrypt/live/api.kwontae.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.kwontae.me/privkey.pem
export MYSQL_ADDR=myuserstore:3306
export MYSQL_DATABASE=myuserstore
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
docker stop myuserstore
docker stop devredis
docker stop myrabbitmq
docker stop gateway
docker stop summaryservice
docker stop messageservice

docker system prune -f

docker network create kwontaeNet

docker pull taehyun123/myuserstore
docker run -d --network kwontaeNet --name myuserstore --restart unless-stopped \
 -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
 -e MYSQL_DATABASE=$MYSQL_DATABASE \
 taehyun123/myuserstore --wait_timeout=31536000

export REDISADDR=devredis:6379
export SESSIONKEY=$(openssl rand -base64 18)
export JAVA_MYSQL_ADDR=myuserstore:3306
export JAVA_MYSQL_DB=myuserstore
export JAVA_MYSQL_PASS=$MYSQL_ROOT_PASSWORD
export JAVA_MYSQL_USER=root

docker run --name devredis -d --restart unless-stopped \
--network kwontaeNet \
redis

docker pull taehyun123/summary
docker run -d --name summaryservice --restart unless-stopped \
--network kwontaeNet \
taehyun123/summary

docker run -d --name myrabbitmq --restart unless-stopped \
--network kwontaeNet \
rabbitmq:3-alpine

docker pull taehyun123/messaging
docker run -d --name messageservice --restart unless-stopped \
-e JAVA_MYSQL_ADDR=$JAVA_MYSQL_ADDR \
-e JAVA_MYSQL_DB=$JAVA_MYSQL_DB \
-e JAVA_MYSQL_PASS=$JAVA_MYSQL_PASS \
-e JAVA_MYSQL_USER=$JAVA_MYSQL_USER \
-e MQHOST=myrabbitmq \
-e MQPORT=5672 \
--network kwontaeNet \
taehyun123/messaging

docker pull taehyun123/gateway
docker run -d -p 443:443 --restart unless-stopped --name gateway \
--network kwontaeNet \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=/etc/letsencrypt/live/api.kwontae.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.kwontae.me/privkey.pem \
-e MYSQL_ADDR=$MYSQL_ADDR \
-e MYSQL_DATABASE=$MYSQL_DATABASE \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
-e SUMMARYADDRS=summaryservice:80 \
-e MESSAGESADDR=messageservice:4000 \
-e MQADDR=myrabbitmq:5672 \
taehyun123/gateway


