# Start by building the application.
FROM golang:1.25.1 AS build

WORKDIR /go/src/app
COPY . .
RUN go get github.com/aritradevelops/billbharat/backend/shared
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/app /
EXPOSE 8080
CMD ["/app"]

