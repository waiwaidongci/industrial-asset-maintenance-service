FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/asset-maintenance ./cmd/server
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/asset-maintenance /asset-maintenance
COPY configs /configs
EXPOSE 8085
ENTRYPOINT ["/asset-maintenance"]
