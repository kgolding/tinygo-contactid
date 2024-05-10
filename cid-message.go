package main

import (
	"errors"
	"strconv"
	"strings"
)

var errInvalidSentence = errors.New("invalid sentence")

type CIDMessage struct {
	Account    int
	EventID    byte
	EventType  byte
	EventCode  CIDEventCode
	AreaID     byte
	UserOrZone int
}

func (m CIDMessage) String() string {
	return "CID Message:\n\r" +
		"\tAccount:    " + strconv.Itoa(m.Account) + "\n\r" +
		"\tEvent ID:   " + strconv.Itoa(int(m.EventID)) + "\n\r" +
		"\tEvent Type: " + strconv.Itoa(int(m.EventType)) + "\n\r" +
		"\tEvent Code: " + m.EventCode.String() + "\n\r" +
		"\tArea ID:    " + strconv.Itoa(int(m.AreaID)) + "\n\r" +
		"\tUser/Zone:  " + strconv.Itoa(int(m.UserOrZone)) + "\n\r"
}

func parseCID(sentence string) (*CIDMessage, error) {
	/*
	   From the documentation:

	   Message Example: AL00123418162701000##(cr)(lf)

	   AL00 = Message start
	   1234 = Account number
	   18 = Contact ID event
	   1 = new event (3 = restore event)
	   627 = Contact ID event code
	   01 = Contact ID event area
	   000 = User or Zone number for contact ID event
	   ##(cr)(lf)= End of message (cr – carriage return 0x13 & lf - line feed 0x0A characters)

	   The response to these messages should be ‘OK’ to acknowledge message received.
	   These are sent as and when they happen, if no OK is received the message is repeated 3 times.

	   Actual observed data started with just "AL":
	   Actual rx data:  AL023418314501000##
	*/

	if !strings.HasPrefix(sentence, "AL") {
		return nil, errInvalidSentence
	}

	// The docs say "AL00" but we actually got "AL", so in the latter case
	// we modify the input to look like the docs say
	switch strings.Index(sentence, "##") {
	case 17: // "AL..."
		sentence = "AL00" + sentence[2:]
	case 19: // "AL00..."
		// Nothing to do
	default:
		return nil, errInvalidSentence
	}

	msg := &CIDMessage{}

	v, err := strconv.Atoi(sentence[4:8])
	if err != nil {
		return nil, errors.New("bad account number")
	}
	msg.Account = v

	v, err = strconv.Atoi(sentence[8:10])
	if err != nil {
		return nil, errors.New("bad event ID")
	}
	msg.EventID = byte(v)

	v, err = strconv.Atoi(sentence[10:11])
	if err != nil {
		return nil, errors.New("bad event type")
	}
	msg.EventType = byte(v)

	v, err = strconv.Atoi(sentence[11:14])
	if err != nil {
		return nil, errors.New("bad event code")
	}
	msg.EventCode = CIDEventCode(v)

	v, err = strconv.Atoi(sentence[14:16])
	if err != nil {
		return nil, errors.New("bad event area")
	}
	msg.AreaID = byte(v)

	v, err = strconv.Atoi(sentence[16:19])
	if err != nil {
		return nil, errors.New("bad user/zone")
	}
	msg.UserOrZone = v

	return msg, nil
}
