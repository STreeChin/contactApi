#-------------------------------
# STEP 1 build executable binary
#-------------------------------
FROM golang:1.14.2-alpine3.11 AS intermediate
MAINTAINER John
LABEL author="John"

WORKDIR /contactapi
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /contactapi/cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /contactapi/contact .

#---------------------------
# STEP 2 build a small image
#---------------------------
FROM scratch
MAINTAINER John
LABEL author="John"

WORKDIR /root/

COPY --from=intermediate /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=intermediate /contactapi/contact .

EXPOSE 8080 8080
CMD ["./contact"]

