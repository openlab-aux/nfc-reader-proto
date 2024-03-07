package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/clausecker/nfc/v2"
	"github.com/skythen/apdu"
)

func main() {
	device, err := nfc.Open("")
	if err != nil {
		log.Fatal("error opening nfc device", err)
	}

	err = device.InitiatorInit()
	if err != nil {
		log.Fatal("error initializing initiator", err)
	}

	modulation := nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}

	targets, err := device.InitiatorListPassiveTargets(modulation)
	if err != nil {
		log.Fatal("Error listing Passive Targets: ", err)
	}
	if len(targets) <= 0 {
		log.Fatal("No targets found")
	}
	target := targets[0].(*nfc.ISO14443aTarget)
	log.Print("found UID: ", target.UID[:target.UIDLen])
	_, err = device.InitiatorSelectPassiveTarget(
		modulation,
		target.UID[:target.UIDLen],
	)
	if err != nil {
		log.Fatal("Error Selecting Target: ", err)
	}

	selectApplication := apdu.Capdu{
		Cla: 0x00,
		Ins: 0xA4,
		P1:  04,
		P2:  00,
		// Data: []byte{0xF0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x07},
		Data: []byte{0xA0, 0x00, 0xDA, 0xDA, 0xDA, 0xDA, 0xDA},
	}

	tx, err := selectApplication.Bytes()
	if err != nil {
		log.Fatal("Error assembling SelectApplication APDU", err)
	}

	rx := make([]byte, 256)

	n, err := device.InitiatorTransceiveBytes(
		tx, rx, 0,
	)
	if err != nil {
		log.Fatal("Error sending bytes: ", err)
	}

	log.Infof("Received %d bytes: %x", n, rx[0:n])
}
