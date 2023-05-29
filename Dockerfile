FROM golang:1.20.4

WORKDIR /boost

COPY go.mod go.mod
COPY go.sum go.sum

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go install -v /boost/cmd/boost && \
    CGO_ENABLED=0 GOOS=linux go install -v /boost/cmd/searcher && \
    rm -rf *

CMD ["boost"]
