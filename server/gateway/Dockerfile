# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server
FROM alpine:latest
RUN apk update && apk add ca-certificates
COPY gateway /gateway
EXPOSE 443
ENTRYPOINT ["/gateway"]