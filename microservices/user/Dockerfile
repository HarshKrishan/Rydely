FROM --platform=linux/amd64 golang:alpine as builder
RUN apk update && apk add --no-cache git
RUN mkdir -p $GOPATH/src/user-service
ADD . $GOPATH/src/user-service
WORKDIR $GOPATH/src/user-service
RUN go get -d -v
RUN go build -o user-service .
# Stage 2
FROM alpine
RUN mkdir /app
COPY --from=builder /go/src/user-service/user-service /app/
COPY --from=builder /go/src/user-service/.env /app/
ARG APP_VERSION
ARG APP_NAME
ARG MODULE_NAME
ENV APP_VERSION=$APP_VERSION
ENV MODULE_NAME = $MODULE_NAME
ENV APP_NAME = $APP_NAME
EXPOSE 8080
WORKDIR /app
CMD ["./user-service"]