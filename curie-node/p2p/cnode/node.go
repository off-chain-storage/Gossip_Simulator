package cnode

import (
	"net"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

type Node struct {
	Addr   net.IP
	Port   int
	PeerID peer.ID
}

// 270.132.13.23:7070@dfdfadjf
func Parse(input string) (*Node, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid address format")
	}

	ip := net.ParseIP(parts[0])
	if ip == nil {
		return nil, errors.New("invalid IP Address")
	}

	portAndID := parts[1]
	subParts := strings.Split(portAndID, "@")
	if len(subParts) != 2 {
		return nil, errors.New("invalid format")
	}

	port, err := strconv.Atoi(subParts[0])
	if err != nil {
		return nil, errors.Wrap(err, "invalid port number")
	}

	id, err := peer.Decode(subParts[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid peer ID")
	}

	return &Node{Addr: ip, Port: port, PeerID: id}, nil
}

func (n *Node) IP() net.IP { return n.Addr }

func (n *Node) TCP() int { return n.Port }

func (n *Node) Peer() peer.ID { return n.PeerID }
