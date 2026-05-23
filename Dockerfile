FROM golang:1.22.1-bookworm AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/booking-bff .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/booking-bff /booking-bff
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/booking-bff"]
