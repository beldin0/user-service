FROM golang:buster AS builder

RUN mkdir /app
WORKDIR /app

# add the dependancies
ADD go.mod go.mod
ADD go.sum go.sum

ADD ./src /app/src

RUN CGO_ENABLED=0 go build -a -o /main /app/src

FROM scratch
COPY --from=builder /main ./
ENTRYPOINT ["./main"]
