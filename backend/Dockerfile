FROM golang:1.26-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE
RUN CGO_ENABLED=0 go build -o /service ./cmd/${SERVICE}

FROM alpine:3.22
COPY --from=build /service /service
ENTRYPOINT ["/service"]
