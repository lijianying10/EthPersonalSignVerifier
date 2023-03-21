FROM 1.20.2-bullseye AS builder
COPY * /app
WORKDIR /app
RUN go build -o /app/main .
FROM ubuntu:22.04
COPY --from=builder /app/main /main
CMD ["/main"]