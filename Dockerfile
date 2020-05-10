FROM golang:buster AS builder

RUN mkdir /app
WORKDIR /app

# add the dependancies
ADD go.mod go.mod
ADD go.sum go.sum

ADD ./src /app/src

RUN CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -a -o /main /app/src

FROM alpine:latest
RUN apk --no-cache add git make build-base
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
