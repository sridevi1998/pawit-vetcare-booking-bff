# PawIt VetCare Booking BFF

Public booking boundary for pet-parent appointment search, slot holds, and booking confirmation. Runs as a Go serverless service on Cloud Run.

## Local Verification

```sh
gofmt -w .
go test ./...
docker build -t pawit-booking-bff:local .
```
