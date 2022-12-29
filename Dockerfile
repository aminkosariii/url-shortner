FROM golang:alpine as builder

RUN mkdir /build

ADD ../../../.config/JetBrains/GoLand2022.3/scratches /build

RUN go get -d -v ./...

RUN go install -v ./...

WORKDIR /build

RUN go build -o main .

#STAGE 2

FROM alpine

RUN adduser -S -D -H /app appuser

USER appuser

COPY ../../../.config/JetBrains/GoLand2022.3/scratches /app

COPY --from-builder /build/main /app/

WORKDIR /app/

EXPOSE 3000

CMD ["./main"]