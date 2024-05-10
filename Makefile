# TARGET=d1mini
TARGET=pico

build:
	tinygo build -target=$(TARGET)

flash:
	tinygo flash -target=$(TARGET)