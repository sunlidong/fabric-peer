/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"fmt"
	"io/ioutil"
	"time"

	"fabricbypeer/common/configtx"
	"fabricbypeer/common/util"
	"fabricbypeer/internal/configtxgen/encoder"
	"fabricbypeer/internal/configtxgen/genesisconfig"
	"fabricbypeer/internal/peer/common"
	"fabricbypeer/internal/pkg/identity"
	"fabricbypeer/protoutil"

	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//ConfigTxFileNotFound channel create configuration tx file not found
type ConfigTxFileNotFound string

func (e ConfigTxFileNotFound) Error() string {
	return fmt.Sprintf("channel create configuration tx file not found %s", string(e))
}

//InvalidCreateTx invalid channel create transaction
type InvalidCreateTx string

func (e InvalidCreateTx) Error() string {
	return fmt.Sprintf("Invalid channel create transaction : %s", string(e))
}

func createCmd(cf *ChannelCmdFactory) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a channel",
		Long:  "Create a channel and write the genesis block to a file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd, args, cf)
		},
	}
	flagList := []string{
		"channelID",
		"file",
		"outputBlock",
		"timeout",
	}
	attachFlags(createCmd, flagList)

	return createCmd
}

func createChannelFromDefaults(cf *ChannelCmdFactory) (*cb.Envelope, error) {
	chCrtEnv, err := encoder.MakeChannelCreationTransaction(
		channelID,
		cf.Signer,
		genesisconfig.Load(genesisconfig.SampleSingleMSPChannelProfile),
	)
	if err != nil {
		return nil, err
	}

	return chCrtEnv, nil
}

func createChannelFromConfigTx(configTxFileName string) (*cb.Envelope, error) {
	cftx, err := ioutil.ReadFile(configTxFileName)
	if err != nil {
		return nil, ConfigTxFileNotFound(err.Error())
	}

	return protoutil.UnmarshalEnvelope(cftx)
}

func sanityCheckAndSignConfigTx(envConfigUpdate *cb.Envelope, signer identity.SignerSerializer) (*cb.Envelope, error) {
	payload, err := protoutil.UnmarshalPayload(envConfigUpdate.Payload)
	if err != nil {
		return nil, InvalidCreateTx("bad payload")
	}

	if payload.Header == nil || payload.Header.ChannelHeader == nil {
		return nil, InvalidCreateTx("bad header")
	}

	ch, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, InvalidCreateTx("could not unmarshall channel header")
	}

	if ch.Type != int32(cb.HeaderType_CONFIG_UPDATE) {
		return nil, InvalidCreateTx("bad type")
	}

	if ch.ChannelId == "" {
		return nil, InvalidCreateTx("empty channel id")
	}

	// Specifying the chainID on the CLI is usually redundant, as a hack, set it
	// here if it has not been set explicitly
	if channelID == "" {
		channelID = ch.ChannelId
	}

	if ch.ChannelId != channelID {
		return nil, InvalidCreateTx(fmt.Sprintf("mismatched channel ID %s != %s", ch.ChannelId, channelID))
	}

	configUpdateEnv, err := configtx.UnmarshalConfigUpdateEnvelope(payload.Data)
	if err != nil {
		return nil, InvalidCreateTx("Bad config update env")
	}

	sigHeader, err := protoutil.NewSignatureHeader(signer)
	if err != nil {
		return nil, err
	}

	configSig := &cb.ConfigSignature{
		SignatureHeader: protoutil.MarshalOrPanic(sigHeader),
	}

	configSig.Signature, err = signer.Sign(util.ConcatenateBytes(configSig.SignatureHeader, configUpdateEnv.ConfigUpdate))
	if err != nil {
		return nil, err
	}

	configUpdateEnv.Signatures = append(configUpdateEnv.Signatures, configSig)

	return protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, signer, configUpdateEnv, 0, 0)
}

//	发送创建通道的Transaction到Order节点
func sendCreateChainTransaction(cf *ChannelCmdFactory) error {
	var err error
	// 定义了一个 Envelope 结构体
	var chCrtEnv *cb.Envelope
	// type Envelope struct {
	// 	#主要就是保存被序列化的有效载荷
	// 	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// 	#由创建者进行的签名信息
	// 	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	// 	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	// 	XXX_unrecognized     []byte   `json:"-"`
	// 	XXX_sizecache        int32    `json:"-"`
	// }

	// 创建 Envelope
	if channelTxFile != "" {
		if chCrtEnv, err = createChannelFromConfigTx(channelTxFile); err != nil {
			return err
		}
	} else {
		// 创建 Envelope
		if chCrtEnv, err = createChannelFromDefaults(cf); err != nil {
			return err
		}
	}
	//校验 刚创建完成的Envelope
	if chCrtEnv, err = sanityCheckAndSignConfigTx(chCrtEnv, cf.Signer); err != nil {
		return err
	}

	// 创建一个用户广播信息的客户端
	var broadcastClient common.BroadcastClient
	broadcastClient, err = cf.BroadcastFactory()
	if err != nil {
		return errors.WithMessage(err, "error getting broadcast client")
	}

	defer broadcastClient.Close()

	//将 创建 通道的Envelope信息广播出去
	err = broadcastClient.Send(chCrtEnv)

	return err
}

func executeCreate(cf *ChannelCmdFactory) error {

	// --1
	err := sendCreateChainTransaction(cf)
	if err != nil {
		return err
	}

	// --2
	block, err := getGenesisBlock(cf)
	if err != nil {
		return err
	}

	// --3
	b, err := proto.Marshal(block)
	if err != nil {
		return err
	}

	file := channelID + ".block"

	// --4
	if outputBlock != common.UndefinedParamValue {
		file = outputBlock
	}

	// --5
	err = ioutil.WriteFile(file, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getGenesisBlock(cf *ChannelCmdFactory) (*cb.Block, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			cf.DeliverClient.Close()
			return nil, errors.New("timeout waiting for channel creation")
		default:
			if block, err := cf.DeliverClient.GetSpecifiedBlock(0); err != nil {
				cf.DeliverClient.Close()
				cf, err = InitCmdFactory(EndorserNotRequired, PeerDeliverNotRequired, OrdererRequired)
				if err != nil {
					return nil, errors.WithMessage(err, "failed connecting")
				}
				time.Sleep(200 * time.Millisecond)
			} else {
				cf.DeliverClient.Close()
				return block, nil
			}
		}
	}
}

func create(cmd *cobra.Command, args []string, cf *ChannelCmdFactory) error {
	// the global chainID filled by the "-c" command
	if channelID == common.UndefinedParamValue {
		return errors.New("must supply channel ID")
	}

	// Parsing of the command line is done so silence cmd usage
	cmd.SilenceUsage = true

	var err error
	if cf == nil {
		cf, err = InitCmdFactory(EndorserNotRequired, PeerDeliverNotRequired, OrdererRequired)
		if err != nil {
			return err
		}
	}
	return executeCreate(cf)
}
