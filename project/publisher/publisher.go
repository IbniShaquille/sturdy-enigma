package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP          = "10.142.0.2"
	serverPort        = "6735"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "lrt-jabodebek-2106706735"
	sslKeyLogFileName = "ssl-key.log"
)


func main() {
	destination := "Harjamukti"
	packetA := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   true,
			IsTrainDeparting:  false,
			TrainNumber:       42,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}

	encodeA := utils.Encoder(packetA)	

	packetB := packetA
	packetB.IsTrainDeparting = true
	packetB.IsTrainArriving = false

	encodeB := utils.Encoder(packetB)


	sslKeyLogFile, err := os.Create(sslKeyLogFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer sslKeyLogFile.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{appLayerProto},
		KeyLogWriter:       sslKeyLogFile,
	}
	connection, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[quic] Dialling from %s to %s\n", connection.LocalAddr(), connection.RemoteAddr())
	fmt.Printf("[quic] Creating receive buffer of size %d\n", bufferSize)
	receiveBuffer := make([]byte, bufferSize)

	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] Opened bidirectional stream %d to %s\n", stream.StreamID(), connection.RemoteAddr())

	fmt.Printf("[quic] [Stream ID: %d] Sending message\n", stream.StreamID())
	_, err = stream.Write([]byte(encodeA))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] [Stream ID: %d] Message sent\n", stream.StreamID())

	stream2, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] Opened bidirectional stream %d to %s\n", stream2.StreamID(), connection.RemoteAddr())

	fmt.Printf("[quic] [Stream ID: %d] Sending message\n", stream2.StreamID())
	_, err = stream2.Write([]byte(encodeB))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] [Stream ID: %d] Message sent\n", stream2.StreamID())


	//cetak packet A
	receiveLength, err := stream.Read(receiveBuffer)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] [Stream ID: %d] Received %d bytes of message from server\n", stream.StreamID(), receiveLength)

	response := receiveBuffer[:receiveLength]
	responseMessage := utils.Decoder(response)
	fmt.Printf("[quic] [Stream ID: %d] Received message: '%+v'\n", stream.StreamID(), responseMessage)

	//cetak packet B
	receiveLength2, err := stream2.Read(receiveBuffer)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] [Stream ID: %d] Received %d bytes of message from server\n", stream2.StreamID(), receiveLength2)

	response2 := receiveBuffer[:receiveLength2]
	responseMessage2 := utils.Decoder(response2)
	fmt.Printf("[quic] [Stream ID: %d] Received message: '%+v'\n", stream2.StreamID(), responseMessage2)


}
