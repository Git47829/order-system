# builder
FROM golang:1.27 as builder
WORKDIR /worker 
COPY go.mod main.go ./
RUN go mod download
RUN go get new
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-api.

#deploy

FROM scratch
COPY --from=builder /worker /worker
USER nonroot:nonroot
EXPOSE 3000 
CMD ["./order-api"]


