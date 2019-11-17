FROM golang:1.13.4 as build

WORKDIR /go/src/github.com/lneves75/krane

COPY . .

RUN CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o krane . && \
    chmod +x krane

FROM alpine

COPY --from=build /go/src/github.com/lneves75/krane/krane /bin/krane

ENTRYPOINT [ "/bin/krane" ]