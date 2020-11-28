FROM alpine:latest AS build
RUN apk upgrade -U -a && \
          apk upgrade && \
          apk add --update go gcc g++ git ca-certificates curl make
WORKDIR /app
COPY ./ /app
RUN printenv
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-X main.version=$RELEASE_VERSION"
