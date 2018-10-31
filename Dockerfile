FROM golang:1.11 as build
WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -tags=jsoniter -o keep4u-backend .

FROM alpine:3.8 AS runtime

COPY --from=build /go/src/keep4u-backend ./
RUN apk add --update ca-certificates

EXPOSE 8080/tcp
ENTRYPOINT ["./keep4u-backend"]
