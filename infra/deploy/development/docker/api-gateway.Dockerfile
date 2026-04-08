FROM alpine
WORKDIR /app

ADD shared shared
ADD build/api-gateway /app/build/api-gateway

ENTRYPOINT ["/app/build/api-gateway"]
