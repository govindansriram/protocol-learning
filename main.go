package main

import "govindansriram/tcip/servers/tcpip"

func main() {
	sett := tcpip.GetDefaultSettings()
	sett.Start()
}
