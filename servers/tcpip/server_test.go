package tcpip

import (
	"net"
	"testing"
	"time"
)

func Test_Start(t *testing.T) {

	signalEnd := make(chan struct{})

	t.Run("test establish connection", func(t *testing.T) {
		go GetDefaultSettings().Start(signalEnd)

		time.Sleep(time.Second * 3) // wait for server to initialize

		_, err := net.Dial("tcp", "localhost:8080")

		if err != nil {
			t.Fatal(err)
		}

		signalEnd <- struct{}{}
		time.Sleep(time.Millisecond * 500) // wait for server to initialize
	})

	t.Run("test send int8", func(t *testing.T) {
		go GetDefaultSettings().Start(signalEnd)

		time.Sleep(time.Second * 1) // wait for server to initialize

		conn, err := net.Dial("tcp", "localhost:8080")

		if err != nil {
			t.Fatal(err)
		}

		bData := []byte{
			0, 0, 1, 2, 4, 7, 9,
		}

		conn.Write(bData)
		time.Sleep(time.Second * 1)

		signalEnd <- struct{}{}
		time.Sleep(time.Millisecond * 500) // wait for server to initialize
	})

	t.Run("send int8 sequences", func(t *testing.T) {
		go GetDefaultSettings().Start(signalEnd)
		time.Sleep(time.Second * 1) // wait for server to initialize

		conn, err := net.Dial("tcp", "localhost:8080")

		if err != nil {
			t.Fatal(err)
		}

		bData := []byte{
			1, 0, 1, 2, 4, 7, 9,
		}

		bData1 := []byte{
			1, 2, 10, 21, 41, 71, 91,
		}

		bData2 := []byte{
			0, 44, 14, 42, 44, 74, 94,
		}

		conn.Write(bData)
		time.Sleep(time.Second * 5)
		conn.Write(bData1)
		time.Sleep(time.Second * 5)
		conn.Write(bData2)

		signalEnd <- struct{}{}
		time.Sleep(time.Millisecond * 500) // wait for server to initialize
	})

	t.Run("send concurrent calls test", func(t *testing.T) {
		go GetDefaultSettings().Start(signalEnd)
		time.Sleep(time.Second * 1) // wait for server to initialize

		sendSignal := func() {
			conn, _ := net.Dial("tcp", "localhost:8080")
			time.Sleep(time.Millisecond * 2000)

			bData := []byte{
				0, 0, 1, 2, 4, 7, 9,
			}
			conn.Write(bData)
		}

		go sendSignal()
		go sendSignal()
		go sendSignal()
		go sendSignal()
		go sendSignal()
		go sendSignal()

		time.Sleep(time.Millisecond * 10_000)

	})
}
