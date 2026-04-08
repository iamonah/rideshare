FROM alpine
WORKDIR /app

ADD shared shared
ADD build/trip-service /app/build/trip-service

ENTRYPOINT ["/app/build/trip-service"]
