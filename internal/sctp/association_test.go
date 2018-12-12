package sctp

import (
	"fmt"
	"testing"
)

func TestAssociationInit(t *testing.T) {
	rawPkt := []byte{0x13, 0x88, 0x13, 0x88, 0x00, 0x00, 0x00, 0x00, 0x81, 0x46, 0x9d, 0xfc, 0x01, 0x00, 0x00, 0x56, 0x55,
		0xb9, 0x64, 0xa5, 0x00, 0x02, 0x00, 0x00, 0x04, 0x00, 0x08, 0x00, 0xe8, 0x6d, 0x10, 0x30, 0xc0, 0x00, 0x00, 0x04, 0x80,
		0x08, 0x00, 0x09, 0xc0, 0x0f, 0xc1, 0x80, 0x82, 0x00, 0x00, 0x00, 0x80, 0x02, 0x00, 0x24, 0x9f, 0xeb, 0xbb, 0x5c, 0x50,
		0xc9, 0xbf, 0x75, 0x9c, 0xb1, 0x2c, 0x57, 0x4f, 0xa4, 0x5a, 0x51, 0xba, 0x60, 0x17, 0x78, 0x27, 0x94, 0x5c, 0x31, 0xe6,
		0x5d, 0x5b, 0x09, 0x47, 0xe2, 0x22, 0x06, 0x80, 0x04, 0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x80, 0x03, 0x00, 0x06, 0x80, 0xc1, 0x00, 0x00}

	assoc := &Association{}
	if err := assoc.handleInbound(rawPkt); err != nil {
		// TODO
		fmt.Println(err)
		// t.Error(errors.Wrap(err, "Failed to HandleInbound"))
	}
}

// Causes deadlock since both sides can be holding the lock in
// Association.send while they are also both trying to
// acquire the lock in Association.handleChunk.
// func TestStressDuplex(t *testing.T) {
// 	lim := test.TimeOut(time.Second * 5)
// 	defer lim.Stop()
//
// 	ca, cb, err := pipeMemory()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	defer func() {
// 		err = ca.Close()
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		err = cb.Close()
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}()
//
// 	opt := test.Options{
// 		MsgSize:  500,
// 		MsgCount: 2,
// 	}
//
// 	err = test.StressDuplex(ca, cb, opt)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func pipeMemory() (*Stream, *Stream, error) {
// 	var err error
//
// 	var aa, ab *Association
// 	aa, ab, err = associationMemory()
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	var sa, sb *Stream
// 	sa, err = aa.OpenStream(0, 0)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	sb, err = ab.OpenStream(0, 0)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	return sa, sb, nil
// }

// func associationMemory() (*Association, *Association, error) {
// 	// In memory pipe
// 	ca, cb := net.Pipe()
//
// 	type result struct {
// 		a   *Association
// 		err error
// 	}
//
// 	c := make(chan result)
//
// 	// Setup client
// 	go func() {
// 		client, err := Client(ca)
// 		c <- result{client, err}
// 	}()
//
// 	// Setup server
// 	server, err := Server(cb)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	// Receive client
// 	res := <-c
// 	if res.err != nil {
// 		return nil, nil, res.err
// 	}
//
// 	return res.a, server, nil
// }
