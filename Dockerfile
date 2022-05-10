FROM golang:1.18 as build
WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o serverd cmd/serverd/main.go
RUN ls -la /app

FROM gcr.io/distroless/base
COPY --from=build /app/serverd /
EXPOSE 8443

CMD ["/serverd", "-tlscert=/etc/certs/tls.crt", "-tlskey=/etc/certs/tls.key", "-port=8443"]