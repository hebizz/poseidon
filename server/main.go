package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "strings"

  "github.com/koron/go-ssdp"
)

var IpList []string

func main() {
  fmt.Println("Scanning for Edge Device...")
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

  m := &ssdp.Monitor{
    Alive:  onAlive,
    Bye:    onBye,
    Search: onSearch,
  }
  if err := m.Start(); err != nil {
    log.Fatal(err)
  }
  select {}
}

func onAlive(m *ssdp.AliveMessage) {
  if m.Type == "edgenode" {
    ip := strings.Split(m.From.String(), ":")
    condition := true
    for _, value := range IpList {
      if value == ip[0] {
        condition = false
        break
      }
    }
    if condition {
      fmt.Printf("IP:%s\n", ip[0])
      IpList = append(IpList, ip[0])
    }
  }
}

func onBye(m *ssdp.ByeMessage) {
  log.Printf("Bye: From=%s Type=%s USN=%s", m.From.String(), m.Type, m.USN)
}

func onSearch(m *ssdp.SearchMessage) {
  //log.Printf("Search: From=%s Type=%s", m.From.String(), m.Type)
}
