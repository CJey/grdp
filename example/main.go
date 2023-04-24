package main

import (
	"fmt"

	"github.com/cjey/grdp"
	"github.com/cjey/grdp/glog"
)

func main() {
	client := grdp.NewClient("192.168.0.3:3889", glog.DEBUG)
	err := client.Login("Administrator", "123456")
	if err != nil {
		fmt.Println("login failed,", err)
	} else {
		fmt.Println("login success")
	}
}
