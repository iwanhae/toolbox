FROM golang:1.20-buster as builder
WORKDIR /app
COPY main.go .
RUN go build -o /main main.go

FROM ubuntu:20.04
RUN apt update 
# Utility
RUN apt install -y vim htop
# Network
RUN apt install -y net-tools iputils-ping dnsutils curl tcpdump iproute2
COPY --from=builder /main /main
EXPOSE 8080
CMD ["/main"]