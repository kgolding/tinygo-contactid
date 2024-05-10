# Contact ID decoder for microcontrollers

A simple proof of concept to use a low power device (tested on ESP8622 & PICO) to decode a security alarms Contact ID data.

## Quick start

1. `git clone git@github.com:kgolding/tinygo-contactid.git`
1. `go mod tidy`
1. Attach PICO (or D1 Mini) via USB
1. `make flash`
1. Open USB serial port in terminal to see log
1. Attach microcontroller to serial output of security alarm