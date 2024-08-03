FROM golang:alpine AS build
RUN apk update && \
    apk add curl \
    git \
    bash \
    make \
    ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./
RUN go mod download
RUN go mod verify

# copy source files and build the binary
COPY . .
RUN make build


FROM alpine:latest
RUN apk --no-cache add ca-certificates bash wget tzdata curl
RUN mkdir -p /var/log/app
WORKDIR /app/
COPY --from=build /app/build .
COPY --from=build /app/build/entrypoint.sh .
COPY --from=build /app/config.yaml .
ENV TZ=Asia/Jakarta
RUN chmod +x ./entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]