package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/clausecker/nfc/v2"
	"github.com/skythen/apdu"
)

var chunk_size = 64

var getMetadataApdu = apdu.Capdu{
	Cla: 0xD0,
	Ins: 0x01,
	P1:  byte(chunk_size),
	P2:  00,
	Ne:  04,
}

var selectApplication = apdu.Capdu{
	Cla:  0x00,
	Ins:  0xA4,
	P1:   04,
	P2:   00,
	Data: []byte{0xA0, 0x00, 0xDA, 0xDA, 0xDA, 0xDA, 0xDA},
}

func transceive(device *nfc.Device, command *apdu.Capdu) (response []byte) {
	tx, err := command.Bytes()
	if err != nil {
		log.Fatal("Error assembling SelectApplication APDU", err)
	}
	log.Infof("tx %d bytes: %x", len(tx), tx)

	rx := make([]byte, 256)

	n, err := device.InitiatorTransceiveBytes(
		tx, rx, 0,
	)
	if err != nil {
		log.Fatal("Error sending bytes: ", err)
	}

	log.Infof("rx %d bytes: %x", n, rx[0:n])

	return rx[0:n]
}

func getToken(device *nfc.Device, len int) (string, error) {

}

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

	transceive(&device, &selectApplication)

	metaResponse := transceive(&device, &getMetadataApdu)

	chunk_count := int(metaResponse[0])
	remainder := int(metaResponse[1])

	log.Infof("Token has %d mod %d chunks = %d bytes", chunk_count, remainder, chunk_count*chunk_size+remainder)

	err = device.Close()
	if err != nil {
		log.Fatal("Error Closing Device: ", err)
	}
}
