# Golang Retome Desktop Protocol

grdp is a pure Golang implementation of the Microsoft RDP (Remote Desktop Protocol) protocol (**client side authorization only**).

## Status

**The project is under development and not finished yet.**

* [ ] SSL Authentication (soon)
* [ ] NLA Authentication

## Example

```golang
client := grdp.NewClient(Host)
err := client.Login(User, Password)
if err != nil {
    fmt.Printf("connect failed: %#v\n", err)
    os.Exit(2)
    return
}
defer client.Close()

fmt.Printf("connected!\n")

sig := make(chan struct{})
once := new(sync.Once)
done := func() {
    once.Do(func() {
        close(sig)
    })
}

client.OnError(func(e error) {
    fmt.Printf("%s Error = %#v\n", time.Now(), e)
    done()
})
client.OnSuccess(func() {
    fmt.Printf("%s Success\n", time.Now())
})
client.OnReady(func() {
    fmt.Printf("%s Ready\n", time.Now())
})
client.OnClose(func() {
    fmt.Printf("%s Close\n", time.Now())
    done()
})
client.OnUpdate(func(_ []pdu.BitmapData) {
    fmt.Printf("%s Update\n", time.Now())
})

fmt.Printf("waiting...\n")
<-sig
```

## Take ideas from

* [rdpy](https://github.com/citronneur/rdpy)
* [node-rdpjs](https://github.com/citronneur/node-rdpjs)
* [gordp](https://github.com/Madnikulin50/gordp)
* [ncrack_rdp](https://github.com/nmap/ncrack/blob/master/modules/ncrack_rdp.cc)
