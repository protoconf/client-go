package mockagent

import (
	"context"
	"log"
	"net"
	"time"

	protoconf_pb "github.com/protoconf/protoconf/agent/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type MockResponse struct {
	wait     time.Duration
	response *protoconf_pb.ConfigUpdate
}

type ProtoconfAgentMock struct {
	protoconf_pb.UnimplementedProtoconfServiceServer
	responses []*MockResponse
	Chan      chan *protoconf_pb.ConfigUpdate
	Stub      protoconf_pb.ProtoconfServiceClient
	initial   proto.Message
}

func NewProtoconfAgentMock(initial proto.Message, responses ...*MockResponse) *ProtoconfAgentMock {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	ctx := context.Background()
	rpcServer := grpc.NewServer()
	srv := &ProtoconfAgentMock{responses: responses, Chan: make(chan *protoconf_pb.ConfigUpdate), initial: initial}
	protoconf_pb.RegisterProtoconfServiceServer(rpcServer, srv)
	go func() {
		context.AfterFunc(ctx, func() {
			rpcServer.GracefulStop()
		})
		if err := rpcServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}
	srv.Stub = protoconf_pb.NewProtoconfServiceClient(conn)
	return srv
}

func (p ProtoconfAgentMock) SubscribeForConfig(request *protoconf_pb.ConfigSubscriptionRequest, srv protoconf_pb.ProtoconfService_SubscribeForConfigServer) error {
	go func() {
		for _, r := range p.responses {
			time.Sleep(r.wait)
			srv.Send(r.response)
		}
	}()
	for {
		select {
		case <-srv.Context().Done():
			return srv.Context().Err()
		case update := <-p.Chan:
			srv.Send(update)
		}
	}
}

func (p ProtoconfAgentMock) SendUpdate(msg proto.Message) {
	v, _ := anypb.New(msg)
	update := &protoconf_pb.ConfigUpdate{
		Value: v,
	}
	go func() {
		p.Chan <- update
	}()
}
