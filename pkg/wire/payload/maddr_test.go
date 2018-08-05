package payload

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/CityOfZion/neo-go/pkg/wire/protocol"
	"github.com/stretchr/testify/assert"
)

func TestNewAddrMessage(t *testing.T) {

	ip := []byte(net.ParseIP("127.0.0.1").To16())

	var ipByte [16]byte
	copy(ipByte[:], ip)

	netaddr, err := NewNetAddr(uint32(time.Now().Unix()), ipByte, 8080, protocol.NodePeerService)
	addrmsg, err := NewAddrMessage()
	addrmsg.AddNetAddr(netaddr)

	buf := new(bytes.Buffer)
	err = addrmsg.EncodePayload(buf)

	addrmsgDec, err := NewAddrMessage()
	r := bytes.NewReader(buf.Bytes())
	err = addrmsgDec.DecodePayload(r)

	assert.Equal(t, nil, err)
	assert.Equal(t, addrmsg.Checksum(), addrmsgDec.Checksum())
}
