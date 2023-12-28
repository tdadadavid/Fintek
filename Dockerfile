FROM golang:1.21.5
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build
CMD [ "make", "start_prod" ]