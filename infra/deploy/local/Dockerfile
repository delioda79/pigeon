FROM golang:1.13 as builder
RUN cd ..
RUN mkdir pigeon
WORKDIR pigeon
COPY . ./
ARG version=dev
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o pigeon ./cmd/pigeon/main.go

FROM scratch
COPY --from=builder /go/pigeon/pigeon .
CMD ["./pigeon"]