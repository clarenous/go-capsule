package node

import (
	"errors"
	"github.com/clarenous/go-capsule/event"
	"github.com/clarenous/go-capsule/p2p"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"

	"github.com/prometheus/prometheus/util/flock"
	log "github.com/sirupsen/logrus"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/clarenous/go-capsule/api"
	cfg "github.com/clarenous/go-capsule/config"
	"github.com/clarenous/go-capsule/consensus"
	"github.com/clarenous/go-capsule/database/leveldb"
	"github.com/clarenous/go-capsule/mining/cpuminer"
	"github.com/clarenous/go-capsule/mining/miningpool"
	"github.com/clarenous/go-capsule/netsync"
	"github.com/clarenous/go-capsule/protocol"
)

const (
	webHost   = "http://127.0.0.1"
	logModule = "node"
)

// Node represent bytom node
type Node struct {
	cmn.BaseService

	config          *cfg.Config
	eventDispatcher *event.Dispatcher
	syncManager     *netsync.SyncManager

	api          *api.API
	chain        *protocol.Chain
	cpuMiner     *cpuminer.CPUMiner
	miningPool   *miningpool.MiningPool
	miningEnable bool
}

// NewNode create bytom node
func NewNode(config *cfg.Config) *Node {
	if err := lockDataDirectory(config); err != nil {
		cmn.Exit("Error: " + err.Error())
	}
	initLogFile(config)
	initActiveNetParams(config)
	initCommonConfig(config)

	// Get store
	if config.DBBackend != "memdb" && config.DBBackend != "leveldb" {
		cmn.Exit(cmn.Fmt("Param db_backend [%v] is invalid, use leveldb or memdb", config.DBBackend))
	}
	coreDB := dbm.NewDB("core", dbm.LevelDBBackend, config.DBDir())
	store := leveldb.NewStore(coreDB)

	dispatcher := event.NewDispatcher()
	txPool := protocol.NewTxPool(store, dispatcher)
	chain, err := protocol.NewChain(store, txPool)
	if err != nil {
		cmn.Exit(cmn.Fmt("Failed to create chain structure: %v", err))
	}

	syncManager, err := netsync.NewSyncManager(config, chain, txPool, dispatcher)
	if err != nil {
		cmn.Exit(cmn.Fmt("Failed to create sync manager: %v", err))
	}

	// run the profile server
	profileHost := config.ProfListenAddress
	if profileHost != "" {
		// Profiling bytomd programs.see (https://blog.golang.org/profiling-go-programs)
		// go tool pprof http://profileHose/debug/pprof/heap
		go func() {
			if err = http.ListenAndServe(profileHost, nil); err != nil {
				cmn.Exit(cmn.Fmt("Failed to register tcp profileHost: %v", err))
			}
		}()
	}

	node := &Node{
		eventDispatcher: dispatcher,
		config:          config,
		syncManager:     syncManager,
		chain:           chain,
		//miningEnable:    config.Mining,
		miningEnable: true,
	}

	node.cpuMiner = cpuminer.NewCPUMiner(chain, txPool, dispatcher)
	node.miningPool = miningpool.NewMiningPool(chain, txPool, dispatcher)

	node.BaseService = *cmn.NewBaseService(nil, "Node", node)

	return node
}

// Lock data directory after daemonization
func lockDataDirectory(config *cfg.Config) error {
	_, _, err := flock.New(filepath.Join(config.RootDir, "LOCK"))
	if err != nil {
		return errors.New("datadir already used by another process")
	}
	return nil
}

func initActiveNetParams(config *cfg.Config) {
	var exist bool
	consensus.ActiveNetParams, exist = consensus.NetParams[config.ChainID]
	if !exist {
		cmn.Exit(cmn.Fmt("chain_id[%v] don't exist", config.ChainID))
	}
}

func initLogFile(config *cfg.Config) {
	if config.LogFile == "" {
		return
	}
	cmn.EnsureDir(filepath.Dir(config.LogFile), 0700)
	file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.WithFields(log.Fields{"module": logModule, "err": err}).Info("using default")
	}

}

func initCommonConfig(config *cfg.Config) {
	cfg.CommonConfig = config
}

// Lanch web broser or not
func launchWebBrowser(port string) {
	webAddress := webHost + ":" + port
	log.Info("Launching System Browser with :", webAddress)
	//if err := browser.Open(webAddress); err != nil {
	//	log.Error(err.Error())
	//	return
	//}
}

func (n *Node) initAndstartAPIServer() error {
	n.api = api.NewAPI(n.chain, n.cpuMiner, n.syncManager)
	return n.api.Start()
}

func (n *Node) OnStart() error {
	if n.miningEnable {
		//if _, err := n.wallet.AccountMgr.GetMiningAddress(); err != nil {
		//	n.miningEnable = false
		//	log.Error(err)
		//} else {
		//	n.cpuMiner.Start()
		//}
		n.cpuMiner.Start()
	}
	if !n.config.VaultMode {
		if err := n.syncManager.Start(); err != nil {
			return err
		}
	}

	err := n.initAndstartAPIServer()
	if err != nil {
		return err
	}

	if !n.config.Web.Closed {
		_, port, err := net.SplitHostPort(n.config.ApiAddress)
		if err != nil {
			log.Error("Invalid api address")
			return err
		}
		launchWebBrowser(port)
	}
	return nil
}

func (n *Node) OnStop() {
	n.BaseService.OnStop()
	if n.miningEnable {
		n.cpuMiner.Stop()
	}
	if !n.config.VaultMode {
		n.syncManager.Stop()
	}
	n.eventDispatcher.Stop()
}

func (n *Node) RunForever() {
	// Sleep forever and then...
	cmn.TrapSignal(func() {
		n.Stop()
	})
}

func (n *Node) NodeInfo() *p2p.NodeInfo {
	return n.syncManager.NodeInfo()
}

func (n *Node) MiningPool() *miningpool.MiningPool {
	return n.miningPool
}
