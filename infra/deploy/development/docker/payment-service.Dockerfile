FROM alpine
WORKDIR /app

ADD shared shared
ADD build/payment-service /app/build/payment-service        

ENTRYPOINT ["/app/build/payment-service"]