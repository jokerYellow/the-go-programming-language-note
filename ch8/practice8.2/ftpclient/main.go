package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

var getfile *os.File
var haveNotWrite bool

func main() {
	address := os.Args[1]
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to :%s", address)

	go func() {
		for {
			response := make([]byte, 10)
			n, e := conn.Read(response)
			if e != nil {
				log.Println(e.Error())
				break
			}
			if getfile != nil {
				if haveNotWrite {
					firstByte := response[0]
					if firstByte != '0' {
						os.Stdout.Write(response[1:n])
						getfile.Close()
						os.Remove(getfile.Name())
						getfile = nil
					}
				} else {
					getfile.Write(response[:n])
				}
			} else {
				os.Stdout.Write(response[:n])
			}
			if n < len(response) && getfile != nil && haveNotWrite == false {
				getfile.Close()
				getfile = nil
			}
			haveNotWrite = false
		}
	}()
	scannerInput := bufio.NewScanner(os.Stdin)
	for {
		if !scannerInput.Scan() {
			return
		}
		b := scannerInput.Bytes()
		cmds := strings.Split(string(b), " ")
		b = append(b, '\n')

		if cmds[0] == "get" {
			f, e := os.Create(cmds[2])
			getfile = f
			haveNotWrite = true
			if e != nil {
				log.Println(e.Error())
				continue
			}
		}
		_, e := conn.Write(b)
		if e != nil {
			log.Printf("%s\n", e)
			return
		}
	}
}
