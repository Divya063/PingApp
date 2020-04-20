package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Divya063/pingApp/utils"
	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	protocolICMP     = 1
	protocolIPv6ICMP = 58
)

type Pinger struct {

	// Count tells pinger to stop after sending (and receiving) Count echo
	// packets. If this option is not specified, pinger will operate until
	// interrupted.
	Count int

	Interval time.Duration

	// Number of packets sent
	PacketsSent int

	// Number of packets received
	PacketsRecv int
	network     string
	// Size of packet being sent
	Size int
	// Source is the source IP address
	Source string

	ipaddr *net.IPAddr
	addr   string

	ipv4 bool
}

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...
const (
	ProtocolICMP     = 1
	ProtocolIPv6ICMP = 58
)

var (
	ipv4Proto = map[string]string{"ip": "ip4:icmp", "udp": "udp4"}
	ipv6Proto = map[string]string{"ip": "ip6:ipv6-icmp", "udp": "udp6"}
)

func NewPinger(addr string) (*Pinger, error) {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}

	var ipv4 bool
	if utils.IsIPv4(ipaddr.IP) {
		ipv4 = true
	} else if utils.IsIPv6(ipaddr.IP) {
		ipv4 = false
	}

	return &Pinger{
		ipaddr:  ipaddr,
		addr:    addr,
		Count:   -1,
		network: "ip",
		ipv4:    ipv4,
	}, nil
}

func listen(listenInterface string, network string) (*icmp.PacketConn, error) {

	c, err := icmp.ListenPacket(network, listenInterface)
	if err != nil {
		return nil, err
	}

	return c, nil

}

func SendAndReceiveRequests(p *Pinger) (*net.IPAddr, time.Duration, float64, error) {
	// Start listening for icmp replies
	loss := 0.0
	var c *icmp.PacketConn
	var err error

	if p.ipv4 {
		c, err = listen("0.0.0.0", ipv4Proto[p.network])
	} else {
		c, err = listen("::", ipv6Proto[p.network])

	}
	if c == nil {
		return nil, 0, 0, err
	}

	defer c.Close()

	// Make a new ICMP message
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, //<< uint(seq), // TODO
			Data: []byte(""),
		},
	}
	msg, err := m.Marshal(nil)
	if err != nil {
		return p.ipaddr, 0, 0, err
	}

	// Send it
	start := time.Now()
	n, err := c.WriteTo(msg, p.ipaddr)
	if err != nil {
		return p.ipaddr, 0, 0, err
	} else if n != len(msg) {
		return p.ipaddr, 0, 0, fmt.Errorf("got %v; want %v", n, len(msg))
	}
	p.PacketsSent++

	var proto int
	if p.ipv4 {
		proto = protocolICMP
	} else {
		proto = protocolIPv6ICMP
	}

	// Wait for a reply
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return p.ipaddr, 0, 0, err
	}
	//n, peer, err := c.ReadFrom(reply)
	if p.ipv4 {

		n, _, _, err = c.IPv4PacketConn().ReadFrom(reply)
	} else {
		n, _, _, err = c.IPv6PacketConn().ReadFrom(reply)
	}
	if err != nil {
		return p.ipaddr, 0, 0, err
	}
	duration := time.Since(start)

	rm, err := icmp.ParseMessage(proto, reply[:n])
	if err != nil {
		return p.ipaddr, 0, 0, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		p.PacketsRecv++
		loss = float64(p.PacketsSent-p.PacketsRecv) / float64(p.PacketsSent) * 100
		return p.ipaddr, duration, loss, nil
	case ipv6.ICMPTypeEchoReply:
		p.PacketsRecv++
		loss = float64(p.PacketsSent-p.PacketsRecv) / float64(p.PacketsSent) * 100
		return p.ipaddr, duration, loss, nil
	default:
		return p.ipaddr, 0, loss, fmt.Errorf("got %+v; want echo reply", rm)
	}
}

func Run(ping *Pinger) {
	counter := 0
	for {
		if ping.Count == counter {
			break
		}
		dst, dur, loss, err := SendAndReceiveRequests(ping)
		if err != nil {
			log.Printf("Ping %s (%s): %s\n", ping.addr, dst, err)
			return
		}
		log.Printf("Ping %s (%s): Duration %s Loss %f\n", ping.addr, dst, dur, loss)

		time.Sleep(ping.Interval * time.Second)

		counter++

	}

}

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:  "ping <host>",
	Long: `The ping command operates by sending Internet Control Message Protocol (ICMP) Echo Request messages to the destination computer and waiting for a response.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Enter a hostname or an IP address")
		}
		count := cmd.Flag("count").Value.String()
		interval := cmd.Flag("interval").Value.String()
		intval, _ := strconv.Atoi(interval)
		addr := args[0]
		ping, _ := NewPinger(addr)

		ping.Count, _ = strconv.Atoi(count)
		ping.Interval = time.Duration(intval)
		Run(ping)

	},
}

func init() {
	pingCmd.PersistentFlags().String("count", "-1", "Count tells pinger to stop after sending (and receiving) Count echo packets. If this option is not specified, pinger will operate until interrupted.")
	pingCmd.PersistentFlags().String("interval", "2", "Periodic interval to send requests")
	RootCmd.AddCommand(pingCmd)

}
