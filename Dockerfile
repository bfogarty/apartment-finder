FROM golang:latest as build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

FROM scratch
COPY --from=build /app /app
USER 1000:1000
ENTRYPOINT ["/app"]
