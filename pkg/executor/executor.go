package executor

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Executor struct {
}

func MakeExecutor() Executor {
	executor := Executor{}
	return executor
}

func (e *Executor) StartClient(destination string) {
	if destination == "" {
		destination = "localhost:50000"
	}
	fmt.Println("Executor connect to: " + destination)

	conn, err := net.Dial("tcp", destination)
	if err != nil {
		log.Error("Could not connect: " + err.Error())
		return
	}
	log.Info("Executor connected")

	// Send initial line
	ex, err := os.Executable()
	if err != nil {
		log.Error("Error: " + err.Error())
		return
	}
	exPath := filepath.Dir(ex)
	pid := strconv.Itoa(os.Getpid())
	line := exPath + ":" + pid + "\n"
	_, err = conn.Write([]byte(line))
	if err != nil {
		log.Error("Error")
		return
	}
	// no answer required

	e.Loop(conn)
}

func (e *Executor) Loop(conn net.Conn) {
	packetExecutor := MakePacketExecutor()

	for {
		// Read
		jsonStr, err := bufio.NewReader(conn).ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("Could not read: " + err.Error())
			break
		}

		fmt.Println("Jsonstr: " + jsonStr)
		packet, err := DecodePacket(jsonStr)
		if err != nil {
			log.Error("Error: ", err.Error())
		}

		// Execute
		err = packetExecutor.Execute(&packet)
		if err != nil {
			log.WithFields(log.Fields{
				"packet": packet,
				"error":  err,
			}).Info("Error executing packet")

			// TODO ERR
		}

		// Answer: Go to JSON
		packetEncoded, err := EncodePacket(packet)
		if err != nil {
			log.Error("Error: ", err.Error())
		}

		n, err := conn.Write(packetEncoded)
		if err != nil {
			log.Error("Error")

			// TODO ERR
		}
		conn.Write([]byte("\n"))
		fmt.Printf("Written: %d bytes", n)

	}
}
