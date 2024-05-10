package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledTimer := time.NewTicker(time.Millisecond * 250)

	serial := machine.UART0
	serial.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.UART_TX_PIN,
		RX:       machine.UART_RX_PIN,
	})

	go func() {
		cid := NewContactIdUart(serial)

		for {
			msg, err := cid.NextMessage()
			if err != nil {
				println(err.Error())
				continue
			}
			println(msg.String())
			led.Set(true)
		}

	}()

	for {
		select {
		case <-ledTimer.C:
			led.Set(!led.Get())
		}
	}
}
