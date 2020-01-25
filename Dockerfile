FROM golang:latest as builder

RUN go get github.com/eclipse/paho.mqtt.golang
RUN go get gopkg.in/yaml.v2
RUN go get github.com/tidwall/gjson

WORKDIR /app
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest
WORKDIR /mqttrepost
COPY --from=builder /app/app ./
ENTRYPOINT ["/mqttrepost/app", "/mqttrepost/config.yml"]
