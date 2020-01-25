# MQTTRepost
## Objective
MQTTRepost subscribes to predefined MQTT topics and reposts all incoming messages in these topics to other predefined topics. You can optionally decode incoming JSON messages so that different values incorporated in the incoming message are posted as separate messages to other topics.

This tool's original use case is to listen to incoming messages from a [zigbee2mqtt](https://www.zigbee2mqtt.io/) gateway and forward them to other MQTT topics for reprocessing. JSON decoding was added to properly restructure incoming data from multisensors (e.g. [Xiaomi Aqara temperature, humidity and pressure sensor](https://www.zigbee2mqtt.io/devices/WSDCGQ11LM.html)) to separate topics that are further processed by [telegraf](https://github.com/influxdata/telegraf) into [influxdb](https://github.com/influxdata/influxdb).

## Building and installing

MQTTRepost is written in Go and will compile to a single standalone binary. It should compile and work both on Linux and on Windows.

For compiling, first install the necessary prerequisites:

    $ go get github.com/eclipse/paho.mqtt.golang
    $ go get github.com/tidwall/gjson
    $ go get gopkg.in/yaml.v2

Then use the following commands to clone the repository and build the binary:

    $ git clone https://github.com/a-bali/mqttrepost.git
    $ cd mqttrepost
    $ go build

This will create the standalone binary named `mqttrepost` that you can place anywhere you like.

### Running with Docker

Docker Hub automatically builds an image from the latest version of MQTTRepost that can be pulled as [`abali/mqttrepost`](https://hub.docker.com/repository/docker/abali/mqttrepost). To use this, map your configuration file to `/mqttrepost/config.yml`:

    $ docker run -v $(pwd)/config.yml:/mqttrepost/config.yml abali/mqttrepost

Alternatively, you can use the supplied Dockerfile to build a container yourself 

    $ git clone https://github.com/a-bali/mqttrepost.git
    $ cd mqttrepost
    $ docker build . -t mqttrepost
    $ docker run -v $(pwd)/config.yml:/mqttrepost/config.yml mqttrepost

## Configuration and usage

For configuration, a YAML formatted file is required. Please use the [sample configuration file](https://raw.githubusercontent.com/a-bali/mqttrepost/master/config.yml) and change it according to your needs. Documentation of various configuration options is included in the sample configuration file.

Once you created a configuration file, MQTTRepost can be launched as follows:

    $ mqttrepost path/to/your/configfile.yml

MQTTRepost will log to standard output. MQTTRepost will not daemonize itself. It is recommended to create a systemd service for it in case you want it running continuously.

## License

MQTTRepost is licensed under GPL 3.0.
