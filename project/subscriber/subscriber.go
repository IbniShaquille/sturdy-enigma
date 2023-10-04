package main

import (
	"fmt"

	"jarkom.cs.ui.ac.id/h01/project/utils"
)

func Handler(packet utils.LRTPIDSPacket) string {
	return ""
}

func main() {
	destination := "Dukuh Atas"
	packet := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   false,
			IsTrainDeparting:  true,
			TrainNumber:       1000,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}
	result := utils.Encoder(packet)
	fmt.Println(result)
	fmt.Println(utils.Decoder(result))
}
