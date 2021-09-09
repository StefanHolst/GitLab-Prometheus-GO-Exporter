# Build
FROM golang:1.16-buster AS build
WORKDIR /build
COPY . .
RUN go build -o /gitlab-prometheus-go-exporter

# Deploy
FROM gcr.io/distroless/base-debian10
COPY --from=build /gitlab-prometheus-go-exporter /gitlab-prometheus-go-exporter
ENTRYPOINT ["/gitlab-prometheus-go-exporter"]