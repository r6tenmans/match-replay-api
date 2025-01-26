# we use stretch linux to build
FROM golang:alpine AS build

# install dependencies
RUN apk update && apk add --no-cache git ca-certificates tzdata upx && update-ca-certificates

# don't use root
ENV USER=r6tm
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /code
# copy dependency files first
# this way docker can skip the go mod download if go.mod or go.sum didn't change
COPY go.mod .
COPY go.sum .
# install go dependencies
RUN go mod download
# copy everything from current working directory into image
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w" ./

# use upx to reduce binary size even further
RUN upx ./backend

# we use "scratch" image to run go service
# the scratch image "doesn't contain anything"
FROM scratch

EXPOSE 8080

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

WORKDIR /app
COPY --from=build /code/backend .
COPY --from=build /code/config/backend.* .
USER r6tm:r6tm

ENTRYPOINT ["/app/backend"]