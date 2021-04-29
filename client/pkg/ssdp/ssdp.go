package ssdp

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/koron/go-ssdp"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
)

//ssdp客户端
func SsdpClient() {
	nt := flag.String("nt", "edgenode", "NT: Type")
	usn := flag.String("usn", "", "USN: ID")
	loc := flag.String("loc", "", "LOCATION: location header")
	srv := flag.String("srv", "", "SERVER:  server header")
	maxAge := flag.Int("maxage", 3600, "cache control, max-age")
	laddr := flag.String("laddr", "", "local address to listen")
	v := flag.Bool("v", false, "verbose mode")
	h := flag.Bool("h", false, "show help")
	flag.Parse()

	if *h {
		flag.Usage()
		return
	}
	if *v {
		ssdp.Logger = log.New(os.Stderr, "[SSDP] ", log.LstdFlags)
	}

	en0, err := net.InterfaceByName(utils.ReadString(config.NetworkCard))
	if err != nil {
		panic(err)
	}
	ssdp.Interfaces = []net.Interface{*en0}

	for {
		err = ssdp.AnnounceAlive(*nt, *usn, *loc, *srv, *maxAge, *laddr)
		if err != nil {
			log.Fatal(err)
		}
		//间隔2s发送数据包
		time.Sleep(time.Duration(2) * time.Second)
	}
}
