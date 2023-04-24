package grdp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/cjey/grdp/core"
	"github.com/cjey/grdp/glog"
	"github.com/cjey/grdp/protocol/nla"
	"github.com/cjey/grdp/protocol/pdu"
	"github.com/cjey/grdp/protocol/sec"
	"github.com/cjey/grdp/protocol/t125"
	"github.com/cjey/grdp/protocol/tpkt"
	"github.com/cjey/grdp/protocol/x224"
)

func init() {
	glog.SetLevel(glog.WARN)
	logger := log.New(os.Stderr, "", 0)
	glog.SetLogger(logger)
}

type Client struct {
	Host string // ip:port
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client

	width  uint16
	height uint16
}

func NewClient(host string) *Client {
	return &Client{
		Host:   host,
		width:  1024,
		height: 768,
	}
}

func split(user string) (domain string, uname string) {
	if strings.Index(user, "\\") != -1 {
		t := strings.Split(user, "\\")
		domain = t[0]
		uname = t[len(t)-1]
	} else if strings.Index(user, "/") != -1 {
		t := strings.Split(user, "/")
		domain = t[0]
		uname = t[len(t)-1]
	} else {
		uname = user
	}
	return
}

func (g *Client) Login(user, pwd string) error {
	conn, err := net.DialTimeout("tcp", g.Host, 3*time.Second)
	if err != nil {
		return errors.New(fmt.Sprintf("[dial err] %v", err))
	}

	domain, user := split(user)

	g.tpkt = tpkt.New(core.NewSocketLayer(conn, nla.NewNTLMv2(domain, user, pwd)))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)

	g.mcs.SetClientDesktop(g.width, g.height)

	g.sec.SetUser(user)
	g.sec.SetPwd(pwd)
	g.sec.SetDomain(domain)

	g.tpkt.SetFastPathListener(g.pdu)
	g.pdu.SetFastPathSender(g.tpkt)

	g.x224.SetRequestedProtocol(x224.PROTOCOL_SSL | x224.PROTOCOL_HYBRID)

	err = g.x224.Connect()
	if err != nil {
		return errors.New(fmt.Sprintf("[x224 connect err] %v", err))
	}
	return nil
}

func (g *Client) OnError(f func(e error)) {
	g.pdu.On("error", f)
}
func (g *Client) OnClose(f func()) {
	g.pdu.On("close", f)
}
func (g *Client) OnSuccess(f func()) {
	g.pdu.On("success", f)
}
func (g *Client) OnReady(f func()) {
	g.pdu.On("ready", f)
}
func (g *Client) OnUpdate(f func([]pdu.BitmapData)) {
	g.pdu.On("update", f)
}

func (g *Client) Close() {
	if g.tpkt != nil {
		g.tpkt.Close()
	}
}
