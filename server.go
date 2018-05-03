package bifrost

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"sync"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/it-chain-Engine/legacy/network/comm/conn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Stream(streamServer pb.StreamService_StreamServer) error {
	//1. RquestPeer를 통해 나에게 Stream연결을 보낸 ConnInfo의정보를 확인
	//2. ConnInfo의정보를정보를 기반으로 Connection을 생성
	//3. 생성완료후 OnConnectionHandler를 통해 처리한다.

	var s struct{}
	envelope, err := bih.createSignedEnvelope(REQUEST_CONNINFO, s)

	err = streamServer.Send(envelope)

	if err != nil {
		return err
	}

	if m, err := recvWithTimeout(3, streamServer); err == nil {

		wg := sync.WaitGroup{}
		wg.Add(1)

		if !IsRequestConnInfoProtocol(m.GetProtocol()) {
			return errors.New(fmt.Sprintf("Not a request connInfo protocol [%s]", m.GetProtocol()))
		}

		log.Printf("Received payload [%s]", envelope.Payload)

		info := bih.info.GetPublicInfo()
		envelope, err := bih.createSignedEnvelope(CONNECTION_ESTABLISH, info)

		if err = streamServer.Send(envelope); err != nil {
			return err
		}

		connectedConnInfo, err := pubConnInfoToConnInfo(envelope.Payload)

		//validate connectedInfo
		if err != nil {
			return err
		}

		_, cf := context.WithCancel(context.Background())
		streamWrapper := stream.NewServerStreamWrapper(streamServer, cf)

		conn, err := conn.NewConnection(*connectedConnInfo, streamWrapper, bih.mux)
		defer conn.Close()

		go func() {
			if err = conn.Start(); err != nil {
				conn.Close()
				wg.Done()
			}
		}()

		bih.onConnectionHandler(conn)

		wg.Wait()
	}

	return nil
}

func Listen(ip string) {

	lis, err := net.Listen("tcp", host.info.Address.IP)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, host)
	reflection.Register(s)

	log.Println("Listen... on: [%s]", host.info.Address.IP)
	func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			s.Stop()
			lis.Close()
		}
	}()
}
