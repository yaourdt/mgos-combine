## build container

FROM golang:1.13.15-alpine AS builder

WORKDIR /go/src/github.com/yaourdt/mgos-combine

COPY . .

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

## production containter

FROM scratch

LABEL author="Mark Dornbach <mark@dornbach.io>"

COPY --from=builder /app ./

EXPOSE 7777

ENTRYPOINT ["./app"]
