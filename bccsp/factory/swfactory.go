/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package factory

import (
	"fabricbypeer/bccsp"
	"fabricbypeer/bccsp/sw"

	"github.com/pkg/errors"
)

const (
	// SoftwareBasedFactoryName is the name of the factory of the software-based BCCSP implementation
	SoftwareBasedFactoryName = "SW"
)

// SWFactory is the factory of the software-based BCCSP.
type SWFactory struct{}

// Name returns the name of this factory
func (f *SWFactory) Name() string {
	return SoftwareBasedFactoryName
}

// Get returns an instance of BCCSP using Opts.
// Get使用Opts返回BCCSP的实例。
func (f *SWFactory) Get(config *FactoryOpts) (bccsp.BCCSP, error) {
	// Validate arguments

	// 判断是否是 nil
	if config == nil || config.SwOpts == nil {
		return nil, errors.New("Invalid config. It must not be nil.")
	}

	// 接参数
	swOpts := config.SwOpts

	// 生成 ks 4934
	var ks bccsp.KeyStore

	// 判断  switch
	switch {

	case swOpts.Ephemeral:
		ks = sw.NewDummyKeyStore()

	case swOpts.FileKeystore != nil:

		fks, err := sw.NewFileBasedKeyStore(nil, swOpts.FileKeystore.KeyStorePath, false)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to initialize software key store")
		}
		ks = fks
	case swOpts.InmemKeystore != nil:
		ks = sw.NewInMemoryKeyStore()
	default:
		// Default to ephemeral key store
		ks = sw.NewDummyKeyStore()
	}

	// 生成 ks 对象， ks 对象是 秘钥储存秘钥 调用  NewWithParams  三个参数： 返回基于软件的BCCSP的一个新实例
	return sw.NewWithParams(swOpts.SecLevel, swOpts.HashFamily, ks)
}

// SwOpts contains options for the SWFactory
type SwOpts struct {
	// Default algorithms when not specified (Deprecated?)
	SecLevel   int    `mapstructure:"security" json:"security" yaml:"Security"`
	HashFamily string `mapstructure:"hash" json:"hash" yaml:"Hash"`

	// Keystore Options
	Ephemeral     bool               `mapstructure:"tempkeys,omitempty" json:"tempkeys,omitempty"`
	FileKeystore  *FileKeystoreOpts  `mapstructure:"filekeystore,omitempty" json:"filekeystore,omitempty" yaml:"FileKeyStore"`
	DummyKeystore *DummyKeystoreOpts `mapstructure:"dummykeystore,omitempty" json:"dummykeystore,omitempty"`
	InmemKeystore *InmemKeystoreOpts `mapstructure:"inmemkeystore,omitempty" json:"inmemkeystore,omitempty"`
}

// Pluggable Keystores, could add JKS, P12, etc..
type FileKeystoreOpts struct {
	KeyStorePath string `mapstructure:"keystore" yaml:"KeyStore"`
}

type DummyKeystoreOpts struct{}

// InmemKeystoreOpts - empty, as there is no config for the in-memory keystore
type InmemKeystoreOpts struct{}
