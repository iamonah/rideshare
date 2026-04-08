FROM alpine
WORKDIR /app

ADD shared shared
ADD build/driver-service /app/build/driver-service

ENTRYPOINT ["/app/build/driver-service"]
