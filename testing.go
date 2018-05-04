package bifrost

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc/metadata"
)

type StreamServer struct {
	countRecv int32
	countSend int32
	peerInfo  PeerInfo
}

func NewMockStreamServer(peerInfo PeerInfo) *StreamServer {
	return &StreamServer{
		countRecv: 0,
		countSend: 0,
		peerInfo:  peerInfo,
	}
}

func (s StreamServer) Send(envelope *pb.Envelope) error {
	log.Print("Mock send func called")

	s.countSend = s.countSend + 1

	if s.countSend == 1 {
		if envelope.Type == pb.Envelope_REQUEST_PEERINFO {
			return nil
		}
		return errors.New("invalid protocol")
	}

	if s.countSend == 2 {
		bool, _, _ := ValidateRequestPeerInfo(envelope)

		if bool {
			return nil
		} else {
			return errors.New("invaild peerinfo")
		}
	}
	return nil
}

func (s StreamServer) Recv() (*pb.Envelope, error) {

	s.countRecv = s.countRecv + 1

	if s.countRecv == 1 {
		payload, _ := json.Marshal(s.peerInfo)

		envelope := &pb.Envelope{}
		envelope.Type = pb.Envelope_REQUEST_PEERINFO
		envelope.Payload = payload
		return envelope, nil
	}

	return nil, nil
}

func (StreamServer) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (StreamServer) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (StreamServer) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (StreamServer) Context() context.Context {
	panic("implement me")
}

func (StreamServer) SendMsg(m interface{}) error {
	panic("implement me")
}

func (StreamServer) RecvMsg(m interface{}) error {
	panic("implement me")
}
