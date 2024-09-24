package main

import (
	"fmt"
	"os"
	"os/signal"
	"proxytun/io_rw"
	"proxytun/tun"
	"strconv"

	"golang.org/x/sys/unix"
)

const (
	ExitSetupFailed = 1
)

const (
	EnvWgTunFd = "WG_TUN_FD"
)

const DefaultMTU = 1420

func main() {
	interfaceName := "tun0"

	tdev, err := func() (tun.Device, error) {
		tunFdStr := os.Getenv(EnvWgTunFd)
		if tunFdStr == "" {
			return tun.CreateTUN(interfaceName, DefaultMTU)
		}

		// construct tun device from supplied fd

		fd, err := strconv.ParseUint(tunFdStr, 10, 32)
		if err != nil {
			return nil, err
		}

		err = unix.SetNonblock(int(fd), true)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}

		file := os.NewFile(uintptr(fd), "")
		return tun.CreateTUNFromFile(file, DefaultMTU)
	}()

	if err == nil {
		realInterfaceName, err2 := tdev.Name()
		if err2 == nil {
			interfaceName = realInterfaceName
		}
	}

	if err != nil {
		fmt.Printf("Failed to create TUN device: %v", err)
		os.Exit(ExitSetupFailed)
	}

	// file := tdev.File()

	// fmt.Println("TUN-DEV file:", file)
	fmt.Println("Everything should be fine...TUN device should exist somewhere")
	
	// fmt.Println("Sleeping for 5 seconds. Please fastly setup tun iface!")
	// time.Sleep(5 * time.Second)

	// Initialize a UDP connection
	// 1. Should create a connection object through with the packets will be sent.
	// 2. Should pass the connection object to the Read routine to send the packets to the network

	go io_rw.RoutineReadFromTun(tdev)

	// for i := 0; i < 20; i++ {
	// 	go io_rw.RoutineWriteToTun(tdev)
	// }


	errs := make(chan error)
	term := make(chan os.Signal, 1)

	signal.Notify(term, unix.SIGTERM)
	// signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	case <-errs:
	}

	fmt.Println("Exiting...")
}
