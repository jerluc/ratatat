package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"github.com/jerluc/ratatat"
	"github.com/jerluc/serial"
)

const (
	DefaultBaudRate = 9600
)

func openSerialPort(devName string, baud int) (io.ReadWriteCloser, error) {
	// TODO: Externalize device baud rate?
	serialCfg := &serial.Config{Name: devName, Baud: 9600}
	serialPort, openErr := serial.OpenPort(serialCfg)
	if openErr != nil {
		return nil, openErr
	}
	return serialPort, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: at-serial <DEVICE> [BAUD_RATE]")
		os.Exit(127)
	}

	devName := os.Args[1]
	baud := DefaultBaudRate
	if len(os.Args) == 3 {
		baud, _ = strconv.Atoi(os.Args[2])
	}

	serialPort, err := openSerialPort(devName, baud)
	if err != nil {
		panic(err)
	}
	defer serialPort.Close()

	commander := ratatat.NewCommander(serialPort)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command (or EXIT to exit): ")
		cmd, _ := reader.ReadString('\n')
		if strings.ToUpper(strings.TrimSpace(cmd)) == "EXIT" {
			break
		}
		res, err := commander.SendAndRecv(cmd)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(res)
		}
	}
}
