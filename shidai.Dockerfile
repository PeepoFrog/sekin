FROM golang:1.22.3 AS shidai-builder

WORKDIR /app

ENV CGO_ENABLED=0 \
		GOOS=linux \
		GOARCH=amd64

COPY ./src/shidai/go.* /app

RUN go mod download

COPY /src/shidai /app

RUN go build -a -tags netgo -installsuffix cgo -o /shidai /app/cmd/main.go

FROM scratch

COPY --from=shidai-builder /shidai /shidai

CMD ["/shidai"]

EXPOSE 8282
