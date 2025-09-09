# Build the backend.
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install sqlc CLI for schemas
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./

# Run sqlc
RUN sqlc generate -f ./server/sqlc.yaml

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./server

# Build the frontend.
FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY client/package*.json ./
RUN npm install
COPY client/ ./
RUN npm run build

# Create the final image.
FROM scratch
COPY --from=builder /server /server
COPY --from=builder /go/bin/goose /goose
COPY --from=frontend-builder /app/dist /client
EXPOSE 8080
CMD [ "/server" ]