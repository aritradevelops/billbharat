# Start by building the application.
FROM golang:1.25.1 AS build

WORKDIR /go/src/app
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12
COPY --from=build /go/src/app/locales /locales
COPY --from=build /go/bin/app /
EXPOSE 8080
CMD ["/app"]

# docker build -t aritradevelops/billbharat-product .
# docker run --env-file .env.docker -p 8080:8080 aritradevelops/billbharat-product
