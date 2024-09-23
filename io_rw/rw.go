package io_rw

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	// "time"

	"proxytun/tun"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)



func RoutineReadFromTun(tdev tun.Device) {
	var (
		// batchSize = 128
		// elemBuf   = make([]byte, (1<<16)-1)
		readErr   error
		// bufs      = make([][]byte, batchSize)
		count     = 0
		// sizes     = make([]int, batchSize)
		offset    = 0
		rxBytes   = 0
		rxCount   = 0

		packet	gopacket.Packet
	)
	ip4, _	:= hex.DecodeString("45")
	ip6, _ 	:= hex.DecodeString("60")

	buffer := make([]byte, (1<<16)-1)

	for {
		// read packets
		count, readErr = tdev.Read(buffer, offset)
		rxCount += 1
		rxBytes += count
		toRead := buffer[10:count]

		ip_proto := toRead[:1]

		if bytes.Equal(ip_proto, ip4) {
			packet = gopacket.NewPacket(toRead, layers.LayerTypeIPv4, gopacket.Default)
		} else if bytes.Equal(ip_proto, ip6) {
			packet = gopacket.NewPacket(toRead, layers.LayerTypeIPv6, gopacket.Default)
		}

		msg := hex.EncodeToString(toRead)

		tdev.EnqueuePaquet(buffer[:count])

		fmt.Println("Read: Count: ", rxCount, "Size:", rxBytes)
		fmt.Println("Read: buff: ", msg) 
		fmt.Println("Read: packet:", packet)
		fmt.Println("Read: Proto: ", hex.EncodeToString(ip_proto))

		if readErr != nil {
			fmt.Println(readErr)
			return
		}
	}
}

const hexIpHead = "00000000000000000000450000c8e9004000011196280a000002effffffa800c076c00b4fac1"

func RoutineWriteToTun(tdev tun.Device) {
	var (
		// batchSize = 128
		// elemBuf   = make([]byte, (1<<16)-1)
		// err   error
		// bufs      = make([][]byte, batchSize)
		count     = 0
		// sizes     = make([]int, batchSize)
		// offset    = 16
		// rxBytes   = 0
		// rxCount   = 0
	)

	payload := strings.Repeat("124312437861faaf29030000010203", 100)

	for {
		// toSend := tdev.DequeuePaquet()
		toSend, err := hex.DecodeString(hexIpHead + payload)
		if err != nil {
			fmt.Println("hexToString conversion problematic...")
			fmt.Println(err)
			return
		}

		count, err = tdev.Write(toSend)

		// toWrite := gopacket.NewPacket(toSend, layers.LayerTypeIPv4, gopacket.Default)

		fmt.Println("Write: Count:", count)
		// fmt.Println("Write: Paquet to Send:", toWrite)

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
