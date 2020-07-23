/*
Copyright Digital Asset Holdings, LLC 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/hyperledger/fabric-protos-go/peer"
	"fabricbypeer/internal/peer/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMissingBlockFile(t *testing.T) {
	defer resetFlags()

	resetFlags()

	cmd := joinCmd(nil)
	AddFlags(cmd)
	args := []string{}
	cmd.SetArgs(args)

	assert.Error(t, cmd.Execute(), "expected join command to fail due to missing blockfilepath")
}

func TestJoin(t *testing.T) {
	defer resetFlags()

	InitMSP()
	resetFlags()

	dir, err := ioutil.TempDir("/tmp", "jointest")
	assert.NoError(t, err, "Could not create the directory %s", dir)
	mockblockfile := filepath.Join(dir, "mockjointest.block")
	err = ioutil.WriteFile(mockblockfile, []byte(""), 0644)
	assert.NoError(t, err, "Could not write to the file %s", mockblockfile)
	defer os.RemoveAll(dir)
	signer, err := common.GetDefaultSigner()
	assert.NoError(t, err, "Get default signer error: %v", err)

	mockResponse := &pb.ProposalResponse{
		Response:    &pb.Response{Status: 200},
		Endorsement: &pb.Endorsement{},
	}

	mockEndorserClient := common.GetMockEndorserClient(mockResponse, nil)

	mockCF := &ChannelCmdFactory{
		EndorserClient:   mockEndorserClient,
		BroadcastFactory: mockBroadcastClientFactory,
		Signer:           signer,
	}

	cmd := joinCmd(mockCF)
	AddFlags(cmd)

	args := []string{"-b", mockblockfile}
	cmd.SetArgs(args)

	assert.NoError(t, cmd.Execute(), "expected join command to succeed")
}

func TestJoinNonExistentBlock(t *testing.T) {
	defer resetFlags()

	InitMSP()
	resetFlags()

	signer, err := common.GetDefaultSigner()
	if err != nil {
		t.Fatalf("Get default signer error: %v", err)
	}

	mockResponse := &pb.ProposalResponse{
		Response:    &pb.Response{Status: 200},
		Endorsement: &pb.Endorsement{},
	}

	mockEndorserClient := common.GetMockEndorserClient(mockResponse, nil)

	mockCF := &ChannelCmdFactory{
		EndorserClient:   mockEndorserClient,
		BroadcastFactory: mockBroadcastClientFactory,
		Signer:           signer,
	}

	cmd := joinCmd(mockCF)

	AddFlags(cmd)

	args := []string{"-b", "mockchain.block"}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.Error(t, err, "expected join command to fail")
	assert.IsType(t, GBFileNotFoundErr(err.Error()), err, "expected error type of GBFileNotFoundErr")
}

func TestBadProposalResponse(t *testing.T) {
	defer resetFlags()

	InitMSP()
	resetFlags()

	mockblockfile := "/tmp/mockjointest.block"
	ioutil.WriteFile(mockblockfile, []byte(""), 0644)
	defer os.Remove(mockblockfile)
	signer, err := common.GetDefaultSigner()
	assert.NoError(t, err, "Get default signer error: %v", err)

	mockResponse := &pb.ProposalResponse{
		Response:    &pb.Response{Status: 500},
		Endorsement: &pb.Endorsement{},
	}

	mockEndorserClient := common.GetMockEndorserClient(mockResponse, nil)

	mockCF := &ChannelCmdFactory{
		EndorserClient:   mockEndorserClient,
		BroadcastFactory: mockBroadcastClientFactory,
		Signer:           signer,
	}

	cmd := joinCmd(mockCF)

	AddFlags(cmd)

	args := []string{"-b", mockblockfile}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.Error(t, err, "expected join command to fail")
	assert.IsType(t, ProposalFailedErr(err.Error()), err, "expected error type of ProposalFailedErr")
}

func TestJoinNilCF(t *testing.T) {
	defer viper.Reset()
	defer resetFlags()

	InitMSP()
	resetFlags()

	dir, err := ioutil.TempDir("/tmp", "jointest")
	assert.NoError(t, err, "Could not create the directory %s", dir)
	mockblockfile := filepath.Join(dir, "mockjointest.block")
	defer os.RemoveAll(dir)
	viper.Set("peer.client.connTimeout", 10*time.Millisecond)
	cmd := joinCmd(nil)
	AddFlags(cmd)
	args := []string{"-b", mockblockfile}
	cmd.SetArgs(args)

	err = cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "endorser client failed to connect to")
}
