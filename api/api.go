package api

import (
	"flag"
	"fmt"
	"github.com/clarenous/go-capsule/mining/cpuminer"
	"github.com/clarenous/go-capsule/netsync"
	"github.com/clarenous/go-capsule/protocol"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"strconv"
)

const (
	defaultGRPCPort   = 8867
	defaultHTTPPort   = 8868
	maxMsgSize        = 1024 * 1024 * 64
	GRPCListenAddress = "127.0.0.1"
)

type API struct {
	server      *grpc.Server
	Chain       *protocol.Chain
	Miner       *cpuminer.CPUMiner
	SyncManager *netsync.SyncManager
}

func NewAPI(chain *protocol.Chain, miner *cpuminer.CPUMiner, syncManager *netsync.SyncManager) *API {
	api := &API{
		Chain:       chain,
		Miner:       miner,
		SyncManager: syncManager,
	}
	api.initServer()
	return api
}

func (a *API) initServer() {
	// set the size for receive Msg
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	}
	a.server = grpc.NewServer(opts...)
	RegisterAPIServiceServer(a.server, a)
	reflection.Register(a.server)
}

func (a *API) runGateway() {
	var run = func() error {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard,
			&runtime.JSONPb{OrigName: true, EmitDefaults: true}))
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize))}
		echoEndpoint := flag.String("echo_endpoint", ":"+strconv.Itoa(defaultGRPCPort), "endpoint of Service")
		err := RegisterAPIServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
		if err != nil {
			return err
		}
		httpPort := fmt.Sprintf("%s%d", ":", defaultHTTPPort)
		return http.ListenAndServe(httpPort, mux)
	}

	go func() {
		log.WithFields(log.Fields{"port_grpc": defaultGRPCPort, "port_http": defaultHTTPPort}).Info("starting gateway on port")
		if err := run(); err != nil {
			log.Error("fail on runGateway", err)
		}
	}()
}

func (a *API) Start() error {
	address := fmt.Sprintf("%s%s%d", GRPCListenAddress, ":", defaultGRPCPort)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	go a.server.Serve(listen)
	a.runGateway()

	return nil
}

func (a *API) Stop() {
	a.server.Stop()
}
