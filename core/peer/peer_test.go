/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hyperledger/fabric-protos-go/common"
	"fabricbypeer/bccsp/sw"
	configtxtest "fabricbypeer/common/configtx/test"
	"fabricbypeer/common/metrics/disabled"
	"fabricbypeer/core/comm"
	"fabricbypeer/core/committer/txvalidator/plugin"
	"fabricbypeer/core/deliverservice"
	validation "fabricbypeer/core/handlers/validation/api"
	"fabricbypeer/core/ledger"
	"fabricbypeer/core/ledger/ledgermgmt"
	"fabricbypeer/core/ledger/ledgermgmt/ledgermgmttest"
	"fabricbypeer/core/ledger/mock"
	ledgermocks "fabricbypeer/core/ledger/mock"
	"fabricbypeer/core/transientstore"
	"fabricbypeer/gossip/gossip"
	gossipmetrics "fabricbypeer/gossip/metrics"
	"fabricbypeer/gossip/service"
	gossipservice "fabricbypeer/gossip/service"
	peergossip "fabricbypeer/internal/peer/gossip"
	"fabricbypeer/internal/peer/gossip/mocks"
	"fabricbypeer/msp/mgmt"
	msptesttools "fabricbypeer/msp/mgmt/testtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	msptesttools.LoadMSPSetupForTesting()
	rc := m.Run()
	os.Exit(rc)
}

func NewTestPeer(t *testing.T) (*Peer, func()) {
	tempdir, err := ioutil.TempDir("", "peer-test")
	require.NoError(t, err, "failed to create temporary directory")

	// Initialize gossip service
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	signer := mgmt.GetLocalSigningIdentityOrPanic(cryptoProvider)

	messageCryptoService := peergossip.NewMCS(&mocks.ChannelPolicyManagerGetter{}, signer, mgmt.NewDeserializersManager(cryptoProvider), cryptoProvider)
	secAdv := peergossip.NewSecurityAdvisor(mgmt.NewDeserializersManager(cryptoProvider))
	defaultSecureDialOpts := func() []grpc.DialOption { return []grpc.DialOption{grpc.WithInsecure()} }
	defaultDeliverClientDialOpts := []grpc.DialOption{grpc.WithBlock()}
	defaultDeliverClientDialOpts = append(
		defaultDeliverClientDialOpts,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(comm.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(comm.MaxSendMsgSize)))
	defaultDeliverClientDialOpts = append(
		defaultDeliverClientDialOpts,
		comm.ClientKeepaliveOptions(comm.DefaultKeepaliveOptions)...,
	)
	gossipConfig, err := gossip.GlobalConfig("localhost:0", nil)

	gossipService, err := gossipservice.New(
		signer,
		gossipmetrics.NewGossipMetrics(&disabled.Provider{}),
		"localhost:0",
		grpc.NewServer(),
		messageCryptoService,
		secAdv,
		defaultSecureDialOpts,
		nil,
		nil,
		gossipConfig,
		&service.ServiceConfig{},
		&deliverservice.DeliverServiceConfig{
			ReConnectBackoffThreshold:   deliverservice.DefaultReConnectBackoffThreshold,
			ReconnectTotalTimeThreshold: deliverservice.DefaultReConnectTotalTimeThreshold,
		},
	)
	require.NoError(t, err, "failed to create gossip service")

	ledgerMgr, err := constructLedgerMgrWithTestDefaults(filepath.Join(tempdir, "ledgersData"))
	require.NoError(t, err, "failed to create ledger manager")

	assert.NoError(t, err)
	transientStoreProvider, err := transientstore.NewStoreProvider(
		filepath.Join(tempdir, "transientstore"),
	)
	assert.NoError(t, err)
	peerInstance := &Peer{
		GossipService:  gossipService,
		StoreProvider:  transientStoreProvider,
		LedgerMgr:      ledgerMgr,
		CryptoProvider: cryptoProvider,
	}

	cleanup := func() {
		ledgerMgr.Close()
		os.RemoveAll(tempdir)
	}
	return peerInstance, cleanup
}

func TestInitialize(t *testing.T) {
	peerInstance, cleanup := NewTestPeer(t)
	defer cleanup()

	peerInstance.Initialize(
		nil,
		plugin.MapBasedMapper(map[string]validation.PluginFactory{}),
		&ledgermocks.DeployedChaincodeInfoProvider{},
		nil,
		nil,
		runtime.NumCPU(),
	)
}

func TestCreateChannel(t *testing.T) {
	peerInstance, cleanup := NewTestPeer(t)
	defer cleanup()

	var initArg string
	peerInstance.Initialize(
		func(cid string) { initArg = cid },
		plugin.MapBasedMapper(map[string]validation.PluginFactory{}),
		&ledgermocks.DeployedChaincodeInfoProvider{},
		nil,
		nil,
		runtime.NumCPU(),
	)

	testChainID := fmt.Sprintf("mytestchainid-%d", rand.Int())
	block, err := configtxtest.MakeGenesisBlock(testChainID)
	if err != nil {
		fmt.Printf("Failed to create a config block, err %s\n", err)
		t.FailNow()
	}

	err = peerInstance.CreateChannel(testChainID, block, &mock.DeployedChaincodeInfoProvider{}, nil, nil)
	if err != nil {
		t.Fatalf("failed to create chain %s", err)
	}

	assert.Equal(t, testChainID, initArg)

	// Correct ledger
	ledger := peerInstance.GetLedger(testChainID)
	if ledger == nil {
		t.Fatalf("failed to get correct ledger")
	}

	// Get config block from ledger
	block, err = ConfigBlockFromLedger(ledger)
	assert.NoError(t, err, "Failed to get config block from ledger")
	assert.NotNil(t, block, "Config block should not be nil")
	assert.Equal(t, uint64(0), block.Header.Number, "config block should have been block 0")

	// Bad ledger
	ledger = peerInstance.GetLedger("BogusChain")
	if ledger != nil {
		t.Fatalf("got a bogus ledger")
	}

	// Correct PolicyManager
	pmgr := peerInstance.GetPolicyManager(testChainID)
	if pmgr == nil {
		t.Fatal("failed to get PolicyManager")
	}

	// Bad PolicyManager
	pmgr = peerInstance.GetPolicyManager("BogusChain")
	if pmgr != nil {
		t.Fatal("got a bogus PolicyManager")
	}

	channels := peerInstance.GetChannelsInfo()
	if len(channels) != 1 {
		t.Fatalf("incorrect number of channels")
	}
}

func TestDeliverSupportManager(t *testing.T) {
	peerInstance, cleanup := NewTestPeer(t)
	defer cleanup()

	manager := &DeliverChainManager{Peer: peerInstance}

	chainSupport := manager.GetChain("fake")
	assert.Nil(t, chainSupport, "chain support should be nil")

	peerInstance.channels = map[string]*Channel{"testchain": {}}
	chainSupport = manager.GetChain("testchain")
	assert.NotNil(t, chainSupport, "chain support should not be nil")
}

func constructLedgerMgrWithTestDefaults(ledgersDataDir string) (*ledgermgmt.LedgerMgr, error) {
	ledgerInitializer := ledgermgmttest.NewInitializer(ledgersDataDir)

	ledgerInitializer.CustomTxProcessors = map[common.HeaderType]ledger.CustomTxProcessor{
		common.HeaderType_CONFIG: &ConfigTxProcessor{},
	}
	ledgerInitializer.Config.HistoryDBConfig = &ledger.HistoryDBConfig{
		Enabled: true,
	}
	return ledgermgmt.NewLedgerMgr(ledgerInitializer), nil
}
