
FROM golang:alpine as builder

RUN adduser -D -g '' appuser
COPY . $GOPATH/src/github.com/joaquinicolas/meliXmen
WORKDIR $GOPATH/src/github.com/joaquinicolas/meliXmen


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app/meliXmen

FROM scratch
WORKDIR /app
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/meliXmen /app/
USER appuser
CMD ls /app/
CMD /app/meliXmen