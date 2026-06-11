#build go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

#copy dependency file
COPY go.mod ./
RUN go mod download

#copy source code
COPY . .
#compile go binary, cgo disabled for security, avoid cahin of trust issue
RUN CGO_ENABLED=0 GOOS=linux go build -o /login-portal .
FROM scratch
WORKDIR /
#copy compiled binary
COPY --from=builder /login-portal /login-portal
ENTRYPOINT ["/login-portal"]