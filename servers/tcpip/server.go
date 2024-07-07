package tcpip

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var PID int

type Server struct {
	maxConnections           uint8
	maxConnectionTimeSeconds uint16
	maxReadWaitTimeSeconds   uint8
	bufferSizeKB             uint32
	port                     uint16
	stream                   bool
}

type packet struct {
	err error
	pos int
}

func GetDefaultSettings() Server {
	return Server{
		maxConnections:           1,
		maxConnectionTimeSeconds: 1000,
		maxReadWaitTimeSeconds:   5,
		bufferSizeKB:             1,
		port:                     8080,
		stream:                   false,
	}
}

func readConnection(
	connection net.Conn,
	buffer *[]byte,
	packetChan chan<- packet,
) {
	n, err := connection.Read(*buffer)
	packetChan <- packet{
		err: err,
		pos: n,
	}
}

func (s Server) Start() {
	server := fmt.Sprintf("localhost:%d", s.port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	connections := make([]net.Conn, 0, s.maxConnections)

	listener, err := net.Listen("tcp", server)

	if err != nil {
		log.Fatalf("could not start server due to: %v", err)
	}

	fmt.Printf("server running on %s \n", server)

	defer listener.Close()

	go func() {
		<-sigChan
		fmt.Println("\n system exiting gracefully")

		for _, c := range connections {
			c.Close()
		}

		listener.Close()
		os.Exit(0)
	}()

	workers := make(chan struct{}, s.maxConnections)

	for {
		workers <- struct{}{}
		conn, err := listener.Accept()

		currentId := PID
		PID++

		if err != nil {
			<-workers
			log.Printf("experienced a error trying to establish a connection: %v \n", err)
		} else {
			fmt.Printf("connection established with process id: %d \n", currentId)
			connections = append(connections, conn)
			go handleConnection(
				conn,
				workers,
				s.maxConnectionTimeSeconds,
				s.maxReadWaitTimeSeconds,
				s.bufferSizeKB,
				currentId)
		}
	}
}

func handleConnection(
	conn net.Conn,
	channel <-chan struct{},
	maxConnectionTimeSeconds uint16,
	maxReadWaitTimeSeconds uint8,
	bufferSizeKB uint32,
	processId int,
) {

	ctx, parentCancelFunc := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*time.Duration(maxConnectionTimeSeconds)))

	defer parentCancelFunc()
	buffer := make([]byte, bufferSizeKB*1024)
	packetChan := make(chan packet)

	go readConnection(conn, &buffer, packetChan)

	var dtype uint8
	proceed := true
	var err error

	localContext, localCancel := context.WithDeadline(
		ctx,
		time.Now().Add(time.Second*time.Duration(maxReadWaitTimeSeconds)),
	)

	fullBuffer := make([]byte, 0, bufferSizeKB*1024*10)

	clearConn := func(
		closeParent,
		closeLocal,
		closeCon bool,
		localCancelFunc context.CancelFunc) {

		if closeParent {
			parentCancelFunc()
		}

		if closeLocal {
			localCancelFunc()
		}

		if closeCon {
			conn.Close()
			<-channel
		}
	}

	select {
	case packet := <-packetChan:
		currentBuffer := make([]byte, packet.pos)
		copy(currentBuffer, buffer)
		proceed, dtype, err = processState(&currentBuffer, true)

		if err != nil {
			log.Printf("received err %v from process id: %d\n", err, processId)
			clearConn(true, true, true, localCancel)
			return
		}

		fullBuffer = append(fullBuffer, currentBuffer[2:]...)

	case <-ctx.Done():
	case <-localContext.Done():
		log.Printf("connection timed out from process id: %d\n", processId)
		clearConn(true, true, true, localCancel)
		return
	}

	localCancel()

	if !proceed {
		clearConn(true, true, true, localCancel)
		processBytes(&fullBuffer, dtype)
		return
	}

	for {
		fmt.Printf("reading from process id: %d\n", processId)
		localContext, localCancel := context.WithDeadline(
			ctx,
			time.Now().Add(time.Second*time.Duration(maxReadWaitTimeSeconds)),
		)

		go readConnection(conn, &buffer, packetChan)

		select {
		case <-ctx.Done():
		case <-localContext.Done():
			fmt.Printf("connection timed out from process id: %d\n", processId)
			clearConn(true, true, true, localCancel)
			return
		case pack := <-packetChan:
			if pack.err != nil {
				log.Printf("process id %d received err %v \n", processId, pack.err)
				clearConn(true, true, true, localCancel)
				return
			} else {
				currentBuffer := make([]byte, pack.pos)
				copy(currentBuffer, buffer)
				proceed, _, err = processState(&currentBuffer, false)

				if err != nil {
					log.Printf("process id %d received err %v \n", processId, pack.err)
					clearConn(true, true, true, localCancel)
				}

				fullBuffer = append(fullBuffer, currentBuffer[1:]...)

				if !proceed {
					clearConn(true, true, true, localCancel)
					processBytes(&fullBuffer, dtype)
					return
				}
			}
		}
		localCancel()
	}
}

func processBytes(pbuffer *[]byte, dtype uint8) {
	switch dtype {
	case uint8(i8):
		fmt.Println(extractInt8(pbuffer))
	case uint8(ui8):
		fmt.Println(extractUint8(pbuffer))
	case uint8(i16):
		fmt.Println(extractInt16(pbuffer))
	case uint8(ui16):
		fmt.Println(extractUint16(pbuffer))
	case uint8(i32):
		fmt.Println(extractInt32(pbuffer))
	case uint8(ui32):
		fmt.Println(extractUint32(pbuffer))
	case uint8(i64):
		fmt.Println(extractInt64(pbuffer))
	case uint8(ui64):
		fmt.Println(extractUint64(pbuffer))
	case uint8(f32):
		fmt.Println(extractFloat32(pbuffer))
	case uint8(f64):
		fmt.Println(extractFloat64(pbuffer))
	}
}
