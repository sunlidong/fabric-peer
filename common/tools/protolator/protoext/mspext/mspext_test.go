/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mspext_test

import (
	"fabricbypeer/common/tools/protolator"
	"fabricbypeer/common/tools/protolator/protoext/mspext"
)

// ensure structs implement expected interfaces
var (
	_ protolator.VariablyOpaqueFieldProto = &mspext.MSPConfig{}
	_ protolator.DecoratedProto           = &mspext.MSPConfig{}

	_ protolator.VariablyOpaqueFieldProto = &mspext.MSPPrincipal{}
	_ protolator.DecoratedProto           = &mspext.MSPPrincipal{}
)
