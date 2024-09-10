FROM golang:1.21 as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o auto-download ./app/main.go


# Use a minimal base image for the final image
FROM alpine
RUN apk --no-cache add tzdata
WORKDIR /app
COPY --from=builder /usr/src/app/auto-download ./auto-download
CMD ["./auto-download"]