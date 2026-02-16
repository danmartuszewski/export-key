FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /bin/export-key ./cmd/export-key

FROM alpine:3.21

RUN apk add --no-cache bash ncurses

ENV TERM=xterm-256color
ENV CLICOLOR_FORCE=1

COPY --from=builder /bin/export-key /usr/local/bin/export-key
COPY example/.env /root/.secrets.env
COPY example/config.yaml /root/.config/export-key/config.yaml

RUN echo 'eval "$(export-key init bash)"' >> /root/.bashrc && \
    echo 'eval "$(export-key init bash)"' >> /root/.profile

WORKDIR /root
CMD ["bash", "-l"]
