package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func pwdCmd(cmd string, conn net.Conn) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(conn, err.Error())
	}
	fmt.Fprintln(conn, path)
}

func cdCmd(cmd string, conn net.Conn) {
	cmds := strings.Split(cmd, " ")
	if len(cmds) <= 1 {
		return
	}
	dir := cmds[1]
	e := os.Chdir(dir)
	if e != nil {
		log.Println(e)
	}
}

func lsCmd(cmd string, conn net.Conn) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(conn, err.Error())
		return
	}
	files, _ := ioutil.ReadDir(path)
	names := []string{}
	for _, f := range files {
		names = append(names, f.Name())
	}
	fmt.Fprintln(conn, strings.Join(names, "\n"))
}

func catCmd(cmd string, conn net.Conn) {
	cmds := strings.Split(cmd, " ")
	if len(cmds) <= 1 {
		return
	}
	path := cmds[1]
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(conn, err.Error())
		return
	}
	_, err = io.Copy(conn, f)
	if err != nil {
		fmt.Fprintln(conn, err.Error())
		return
	}
}

func getCmd(cmd string, conn net.Conn) {
	cmds := strings.Split(cmd, " ")
	if len(cmds) <= 1 {
		return
	}
	path := cmds[1]
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(conn, "1"+err.Error())
		return
	}
	fmt.Fprintln(conn, "0")
	_, err = io.Copy(conn, f)
	if err != nil {
		fmt.Fprintln(conn, "1"+err.Error())
		return
	}
}

func defaultCmd(cmd string, conn net.Conn) {
	fmt.Fprintf(conn, "invalid command:%s\n", cmd)
}
