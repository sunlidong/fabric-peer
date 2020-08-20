/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"time"

	"fabricbypeer/gossip/election"
	"fabricbypeer/gossip/util"

	"github.com/spf13/viper"
)

const (
	btlPullMarginDefault           = 10
	transientBlockRetentionDefault = 1000
)

// ServiceConfig is the config struct for gossip services

//  服务的配置
type ServiceConfig struct {
	// PeerTLSEnabled enables/disables Peer TLS.

	// peer 节点 是否开启 tls
	PeerTLSEnabled bool
	// Endpoint which overrides the endpoint the peer publishes to peers in its organization.
	//端点，它覆盖对等发布给其组织中的对等的端点。
	Endpoint              string
	NonBlockingCommitMode bool
	// UseLeaderElection defines whenever peer will initialize dynamic algorithm for "leader" selection.
	//使用Leader Election定义了何时peer将初始化“Leader”选择的动态算法。
	UseLeaderElection bool
	// OrgLeader statically defines peer to be an organization "leader".
	// Org Leader静态定义peer为组织“Leader”。
	OrgLeader bool
	// ElectionStartupGracePeriod is the longest time peer waits for stable membership during leader
	// election startup (unit: second).
	//选举启动期是领导期间同僚等待稳定成员的最长时间
	//选举启动(单位:第二)。
	ElectionStartupGracePeriod time.Duration
	// ElectionMembershipSampleInterval is the time interval for gossip membership samples to check its stability (unit: second).
	//选举成员ship样本时间间隔是八卦成员样本检验其稳定性的时间间隔(单位:秒)。
	ElectionMembershipSampleInterval time.Duration
	// ElectionLeaderAliveThreshold is the time passes since last declaration message before peer decides to
	// perform leader election (unit: second).
	//选举领袖在世的门槛是指时间从上次宣布讯息到peer决定之前
	//进行领导人选举(单位:第二)。
	ElectionLeaderAliveThreshold time.Duration
	// ElectionLeaderElectionDuration is the time passes since last declaration message before peer decides to perform
	// leader election (unit: second).
	// 是peer决定执行之前自上次声明消息起经过的时间
	//领导人选举(单位:二)。
	ElectionLeaderElectionDuration time.Duration
	// 确定私有数据对应的最大持续时间
	//给定块。
	// PvtDataPullRetryThreshold determines the maximum duration of time private data corresponding for
	// a given block.
	PvtDataPullRetryThreshold time.Duration
	// PvtDataPushAckTimeout is the maximum time to wait for the acknoledgement from each peer at private
	// data push at endorsement time.
	// 是私有节点等待应答的最大时间
	//背书时数据推送。
	PvtDataPushAckTimeout time.Duration
	// BtlPullMargin is the block to live pulling margin, used as a buffer to prevent peer from trying to pull private data
	// from peers that is soon to be purged in next N blocks.
	// 用作缓冲区防止peer试图提取私有数据
	//来自即将在接下来的N个块中被清除的对等节点。
	BtlPullMargin uint64
	// TransientstoreMaxBlockRetention defines the maximum difference between the current ledger's height upon commit,
	// and the private data residing inside the transient store that is guaranteed not to be purged.
	// 定义了提交时当前分类帐高度的最大差异，
	//和驻留在保证不会被清除的临时存储中的私有数据。
	TransientstoreMaxBlockRetention uint64
	// SkipPullingInvalidTransactionsDuringCommit is a flag that indicates whether pulling of invalid
	// transaction's private data from other peers need to be skipped during the commit time and pulled
	// only through reconciler.
	// 是一个标志，指示拉取是否无效
	//事务来自其他对等点的私有数据需要在提交期间跳过并提取
	//只有通过和解者。
	SkipPullingInvalidTransactionsDuringCommit bool
}

func GlobalConfig() *ServiceConfig {
	c := &ServiceConfig{}
	c.loadGossipConfig()
	return c
}

func (c *ServiceConfig) loadGossipConfig() {

	c.PeerTLSEnabled = viper.GetBool("peer.tls.enabled")
	c.Endpoint = viper.GetString("peer.gossip.endpoint")
	c.NonBlockingCommitMode = viper.GetBool("peer.gossip.nonBlockingCommitMode")
	c.UseLeaderElection = viper.GetBool("peer.gossip.useLeaderElection")
	c.OrgLeader = viper.GetBool("peer.gossip.orgLeader")

	c.ElectionStartupGracePeriod = util.GetDurationOrDefault("peer.gossip.election.startupGracePeriod", election.DefStartupGracePeriod)
	c.ElectionMembershipSampleInterval = util.GetDurationOrDefault("peer.gossip.election.membershipSampleInterval", election.DefMembershipSampleInterval)
	c.ElectionLeaderAliveThreshold = util.GetDurationOrDefault("peer.gossip.election.leaderAliveThreshold", election.DefLeaderAliveThreshold)
	c.ElectionLeaderElectionDuration = util.GetDurationOrDefault("peer.gossip.election.leaderElectionDuration", election.DefLeaderElectionDuration)

	c.PvtDataPushAckTimeout = viper.GetDuration("peer.gossip.pvtData.pushAckTimeout")
	c.PvtDataPullRetryThreshold = viper.GetDuration("peer.gossip.pvtData.pullRetryThreshold")
	c.SkipPullingInvalidTransactionsDuringCommit = viper.GetBool("peer.gossip.pvtData.skipPullingInvalidTransactionsDuringCommit")

	c.BtlPullMargin = btlPullMarginDefault
	if viper.IsSet("peer.gossip.pvtData.btlPullMargin") {
		btlMarginVal := viper.GetInt("peer.gossip.pvtData.btlPullMargin")
		if btlMarginVal >= 0 {
			c.BtlPullMargin = uint64(btlMarginVal)
		}
	}

	c.TransientstoreMaxBlockRetention = uint64(viper.GetInt("peer.gossip.pvtData.transientstoreMaxBlockRetention"))
	if c.TransientstoreMaxBlockRetention == 0 {
		logger.Warning("Configuration key peer.gossip.pvtData.transientstoreMaxBlockRetention isn't set, defaulting to", transientBlockRetentionDefault)
		c.TransientstoreMaxBlockRetention = transientBlockRetentionDefault
	}
}
