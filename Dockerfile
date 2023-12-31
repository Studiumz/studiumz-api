# ----------- STEP 1: Build executable -----------
FROM golang:alpine AS builder

# Install git
RUN apk update && apk add --no-cache git tzdata ca-certificates && update-ca-certificates

# Create unprivileged user
ENV USER=studiumz
ENV UID=1001

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

WORKDIR /app/studiumz-api
COPY . .

# Fetch dependencies
RUN go get -d -v
# Add base env file for mounting config
RUN touch /tmp/.env
# Build executable
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/main main.go

# ----------- STEP 2: Build small image ----------- 
FROM scratch
WORKDIR /app/studiumz-api

# Import user files, certificates, and timezone info
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Import env file config, log file, and executable
COPY --from=builder --chown=studiumz:studiumz /tmp/.env /app/studiumz-api/.env
COPY --from=builder --chown=studiumz:studiumz /app/studiumz-api/bin/main /app/studiumz-api/bin/main
# Change user to unprivileged and set timezone
USER studiumz:studiumz
ENV TZ=Asia/Jakarta

# Run executable
ENTRYPOINT [ "./bin/main" ]
