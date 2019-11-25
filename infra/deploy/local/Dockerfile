FROM golang:1.13 as builder
RUN cd ..
RUN mkdir peristeri
WORKDIR peristeri
COPY . ./
ARG version=dev
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-X main.version=$version" -o peristeri ./cmd/peristeri/main.go

FROM scratch
COPY --from=builder /go/peristeri/peristeri .
CMD ["./peristeri"]