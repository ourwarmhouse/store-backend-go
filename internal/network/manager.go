package network

import (
	"fmt"
	"sync"

	"github.com/pions/pkg/stun"
	"github.com/pions/webrtc/internal/datachannel"
	"github.com/pions/webrtc/internal/dtls"
	"github.com/pions/webrtc/internal/sctp"
	"github.com/pions/webrtc/internal/srtp"
	"github.com/pions/webrtc/pkg/ice"
	"github.com/pions/webrtc/pkg/rtp"
	"github.com/pkg/errors"
)

// Manager contains all network state (DTLS, SRTP) that is shared between ports
// It is also used to perform operations that involve multiple ports
type Manager struct {
	icePwd      []byte
	iceNotifier ICENotifier

	dtlsState *dtls.State

	certPairLock sync.RWMutex
	certPair     *dtls.CertPair

	bufferTransportGenerator BufferTransportGenerator
	bufferTransports         map[uint32]chan<- *rtp.Packet

	// https://tools.ietf.org/html/rfc3711#section-3.2.3
	// A cryptographic context SHALL be uniquely identified by the triplet
	//  <SSRC, destination network address, destination transport port number>
	// contexts are keyed by IP:PORT:SSRC
	srtpContextsLock sync.RWMutex
	srtpContexts     map[string]*srtp.Context

	sctpAssociation *sctp.Association

	ports []*port
}

// NewManager creates a new network.Manager
func NewManager(icePwd []byte, bufferTransportGenerator BufferTransportGenerator, dataChannelEventHandler DataChannelEventHandler) (m *Manager, err error) {
	m = &Manager{
		bufferTransports:         make(map[uint32]chan<- *rtp.Packet),
		srtpContexts:             make(map[string]*srtp.Context),
		bufferTransportGenerator: bufferTransportGenerator,
		icePwd: icePwd,
	}
	m.dtlsState, err = dtls.NewState(true)
	if err != nil {
		return nil, err
	}

	m.sctpAssociation = sctp.NewAssocation(func(raw []byte) {
		for _, p := range m.ports {
			if p.iceState == ice.ConnectionStateCompleted {
				p.sendSCTP(raw)
				return
			}
		}
	}, func(data []byte, streamIdentifier uint16) {
		msg, err := datachannel.Parse(data)
		if err != nil {
			fmt.Println(errors.Wrap(err, "Failed to parse DataChannel packet"))
			return
		}
		switch m := msg.(type) {
		case *datachannel.ChannelOpen:
			dataChannelEventHandler(&DataChannelCreated{streamIdentifier: streamIdentifier, Label: string(m.Label)})
		case *datachannel.Data:
			dataChannelEventHandler(&DataChannelMessage{streamIdentifier: streamIdentifier, Body: m.Data})
		default:
			fmt.Println("Unhandled DataChannel message", m)
		}
	})

	return m, err
}

// Listen starts a new Port for this manager
func (m *Manager) Listen(address string) (boundAddress *stun.TransportAddr, err error) {
	p, err := newPort(address, m)
	if err != nil {
		return nil, err
	}

	m.ports = append(m.ports, p)
	return p.listeningAddr, nil
}

// Close cleans up all the allocated state
func (m *Manager) Close() {
	m.sctpAssociation.Close()
	m.dtlsState.Close()
}

// DTLSFingerprint generates the fingerprint included in an SessionDescription
func (m *Manager) DTLSFingerprint() string {
	return m.dtlsState.Fingerprint()
}

// SendRTP finds a connected port and sends the passed RTP packet
func (m *Manager) SendRTP(packet *rtp.Packet) {
	for _, p := range m.ports {
		if p.iceState == ice.ConnectionStateCompleted {
			p.sendRTP(packet)
			return
		}
	}
}

func (m *Manager) iceHandler(p *port) {
}
