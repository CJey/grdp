package main

import (
	"fmt"
	"time"

	"github.com/cjey/grdp"
	"github.com/cjey/grdp/protocol/pdu"
)

func main() {
	defer func() {
		fmt.Printf("panic = %#v\n", recover())
	}()

	client := grdp.NewClient("10.234.74.143:3389")
	err := client.Login("Administrator", "CJey.Hou0723")
	if err != nil {
		fmt.Printf("login failed: %#v\n", err)
		return
	}

	fmt.Printf("connected!\n")

	done := make(chan struct{})

	client.OnError(func(e error) {
		fmt.Printf("%s onError = %#v\n", time.Now(), e)
	})
	client.OnSuccess(func() {
		fmt.Printf("%s onSuccess\n", time.Now())
	})
	client.OnReady(func() {
		fmt.Printf("%s onReady\n", time.Now())
	})
	client.OnClose(func() {
		fmt.Printf("%s onClose\n", time.Now())
	})
	client.OnUpdate(func(_ []pdu.BitmapData) {
		fmt.Printf("%s onUpdate\n", time.Now())
	})

	fmt.Printf("waiting...\n")
	<-done
}
