[![Build Status](https://travis-ci.org/kost/smtp2mqtt.png)](https://travis-ci.org/kost/smtp2mqtt)
[![Circle Status](https://circleci.com/gh/kost/smtp2mqtt.svg?style=shield&circle-token=:circle-token)](https://circleci.com/gh/kost/smtp2mqtt)

# smtp2mqtt

Simple SMTP to MQTT relay/forwarder.

Ever wanted to just forward your alert e-mails to MQTT? smtp2mqtt can do it.

# Features

-   Single executable (thanks to Go!)
-   Linux/Windows/Mac/BSD support
-   multipart text/html e-mails supported
-   (optional) JSON support

# Examples

Here is quick example, just to get idea what you can do with it.

## Quick Examples

Start server:

    $ ./smtp2mqtt -listen 0.0.0.0:10025
    2019/10/22 21:51:25 Listening on: 0.0.0.0:10025

Run SMTP server forwarder and forward to MQTT 192.168.1.1 in JSON:

    $ ./smtp2mqtt -json -topic smtp/ -mqtt tcp://192.168.1.1:1883
    2019/10/22 21:52:25 Listening on: 0.0.0.0:10025

# Download

You can find binary and source releases on Github under "Releases". Here's the [link to the latest release](https://github.com/kost/smtp2mqtt/releases/latest)

# Options explained

    $ ./smtp2mqtt
    Usage of ./smtp2mqtt:
      -allow string
        	Allow only specific IPs to send e-mail (e.g. 192.168.1.)
      -debug
        	Enable debug messages
      -deny string
        	Deny specific IPs to send e-mail (e.g. 192.168.1.10)
      -json
        	post to MQTT topic as json
      -keep
        	keep connection to MQTT
      -listen string
        	Listen on specific IP and port (default "0.0.0.0:10025")
      -mqtt string
        	connect to specified MQTT server (default "tcp://127.0.0.1:1883")
      -password string
        	MQTT password for connecting
      -topic string
        	prepend specified string to MQTT topic (e.g. 'smtp/')
      -user string
        	MQTT username for connecting
      -welcome string
        	Welcome message for SMTP session (default "MQTT-forwarder ESMTP ready.")

# Building

## Linux/Mac/POSIX builds

Just type:

    go build

Static compiling:

    CGO_ENABLED=0 go build -ldflags "-extldflags -static"

## Windows builds

Just type:

    go build

### ToDo

-   [ ] Implement TLS support for MQTT

### Done

-   [x] Implement JSON support

# Credits

Vlatko Kosturjak
