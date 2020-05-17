FROM golang:buster AS builder

RUN mkdir /app
WORKDIR /app

# add the dependancies
COPY go.mod go.mod
COPY go.sum go.sum

COPY ./src /app/src

RUN CGO_ENABLED=0 go build -a -o /main /app/src

FROM scratch
COPY --from=builder /main ./
ENTRYPOINT ["./main"]
