// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

package aead

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/tsingson/tink/golang/core/registry"
	"github.com/tsingson/tink/golang/keyset"
	"github.com/tsingson/tink/golang/subtle/aead"
	"github.com/tsingson/tink/golang/subtle/mac"
	"github.com/tsingson/tink/golang/subtle/random"
	ctrpb "github.com/tsingson/tink/proto/aes_ctr_go_proto"
	aeadpb "github.com/tsingson/tink/proto/aes_ctr_hmac_aead_go_proto"
	commonpb "github.com/tsingson/tink/proto/common_go_proto"
	hmacpb "github.com/tsingson/tink/proto/hmac_go_proto"
	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

const (
	aesCTRHMACAEADKeyVersion = 0
	aesCTRHMACAEADTypeURL    = "type.googleapis.com/google.crypto.tink.AesCtrHmacAeadKey"
	minHMACKeySizeInBytes    = 16
	minTagSizeInBytes        = 10
)

// common errors
var errInvalidAESCTRHMACAEADKey = fmt.Errorf("aes_ctr_hmac_aead_key_manager: invalid key")
var errInvalidAESCTRHMACAEADKeyFormat = fmt.Errorf("aes_ctr_hmac_aead_key_manager: invalid key format")

// aesCTRHMACAEADKeyManager is an implementation of KeyManager interface.
// It generates new AESCTRHMACAEADKey keys and produces new instances of EncryptThenAuthenticate subtle.
type aesCTRHMACAEADKeyManager struct{}

// Assert that aesCTRHMACAEADKeyManager implements the KeyManager interface.
var _ registry.KeyManager = (*aesCTRHMACAEADKeyManager)(nil)

// newAESCTRHMACAEADKeyManager creates a new aesCTRHMACAEADKeyManager.
func newAESCTRHMACAEADKeyManager() *aesCTRHMACAEADKeyManager {
	return new(aesCTRHMACAEADKeyManager)
}

// Primitive creates an AEAD for the given serialized AESCTRHMACAEADKey proto.
func (km *aesCTRHMACAEADKeyManager) Primitive(serializedKey []byte) (interface{}, error) {
	if len(serializedKey) == 0 {
		return nil, errInvalidAESCTRHMACAEADKey
	}
	key := new(aeadpb.AesCtrHmacAeadKey)
	if err := proto.Unmarshal(serializedKey, key); err != nil {
		return nil, errInvalidAESCTRHMACAEADKey
	}
	if err := km.validateKey(key); err != nil {
		return nil, err
	}

	ctr, err := aead.NewAESCTR(key.AesCtrKey.KeyValue, int(key.AesCtrKey.Params.IvSize))
	if err != nil {
		return nil, fmt.Errorf("aes_ctr_hmac_aead_key_manager: cannot create new primitive: %v", err)
	}

	hmacKey := key.HmacKey
	hmac, err := mac.NewHMAC(hmacKey.Params.Hash.String(), hmacKey.KeyValue, hmacKey.Params.TagSize)
	if err != nil {
		return nil, fmt.Errorf("aes_ctr_hmac_aead_key_manager: cannot create mac primitive, error: %v", err)
	}

	aead, err := aead.NewEncryptThenAuthenticate(ctr, hmac, int(hmacKey.Params.TagSize))
	if err != nil {
		return nil, fmt.Errorf("aes_ctr_hmac_aead_key_manager: cannot create encrypt then authenticate primitive, error: %v", err)
	}
	return aead, nil
}

// NewKey creates a new key according to the given serialized AesCtrHmacAeadKeyFormat.
func (km *aesCTRHMACAEADKeyManager) NewKey(serializedKeyFormat []byte) (proto.Message, error) {
	if len(serializedKeyFormat) == 0 {
		return nil, errInvalidAESCTRHMACAEADKeyFormat
	}
	keyFormat := new(aeadpb.AesCtrHmacAeadKeyFormat)
	if err := proto.Unmarshal(serializedKeyFormat, keyFormat); err != nil {
		return nil, errInvalidAESCTRHMACAEADKeyFormat
	}
	if err := km.validateKeyFormat(keyFormat); err != nil {
		return nil, fmt.Errorf("aes_ctr_hmac_aead_key_manager: invalid key format: %v", err)
	}
	return &aeadpb.AesCtrHmacAeadKey{
		Version: aesCTRHMACAEADKeyVersion,
		AesCtrKey: &ctrpb.AesCtrKey{
			Version:  aesCTRHMACAEADKeyVersion,
			KeyValue: random.GetRandomBytes(keyFormat.AesCtrKeyFormat.KeySize),
			Params:   keyFormat.AesCtrKeyFormat.Params,
		},
		HmacKey: &hmacpb.HmacKey{
			Version:  aesCTRHMACAEADKeyVersion,
			KeyValue: random.GetRandomBytes(keyFormat.HmacKeyFormat.KeySize),
			Params:   keyFormat.HmacKeyFormat.Params,
		},
	}, nil
}

// NewKeyData creates a new KeyData according to specification in the given serialized
// AesCtrHmacAeadKeyFormat.
// It should be used solely by the key management API.
func (km *aesCTRHMACAEADKeyManager) NewKeyData(serializedKeyFormat []byte) (*tinkpb.KeyData, error) {
	key, err := km.NewKey(serializedKeyFormat)
	if err != nil {
		return nil, err
	}
	serializedKey, err := proto.Marshal(key)
	if err != nil {
		return nil, err
	}
	return &tinkpb.KeyData{
		TypeUrl:         km.TypeURL(),
		Value:           serializedKey,
		KeyMaterialType: tinkpb.KeyData_SYMMETRIC,
	}, nil
}

// DoesSupport indicates if this key manager supports the given key type.
func (km *aesCTRHMACAEADKeyManager) DoesSupport(typeURL string) bool {
	return typeURL == aesCTRHMACAEADTypeURL
}

// TypeURL returns the key type of keys managed by this key manager.
func (km *aesCTRHMACAEADKeyManager) TypeURL() string {
	return aesCTRHMACAEADTypeURL
}

// validateKey validates the given AesCtrHmacAeadKey proto.
func (km *aesCTRHMACAEADKeyManager) validateKey(key *aeadpb.AesCtrHmacAeadKey) error {
	err := keyset.ValidateKeyVersion(key.Version, aesCTRHMACAEADKeyVersion)
	if err != nil {
		return fmt.Errorf("aes_ctr_hmac_aead_key_manager: %v", err)
	}

	// Validate AesCtrKey.
	keySize := uint32(len(key.AesCtrKey.KeyValue))
	if err := aead.ValidateAESKeySize(keySize); err != nil {
		return fmt.Errorf("aes_ctr_hmac_aead_key_manager: %v", err)
	}
	params := key.AesCtrKey.Params
	if params.IvSize < aead.AESCTRMinIVSize || params.IvSize > 16 {
		return errors.New("aes_ctr_hmac_aead_key_manager: invalid AesCtrHmacAeadKey: IV size out of range")
	}
	return nil
}

// validateKeyFormat validates the given AesCtrHmacAeadKeyFormat proto.
func (km *aesCTRHMACAEADKeyManager) validateKeyFormat(format *aeadpb.AesCtrHmacAeadKeyFormat) error {
	// Validate AesCtrKeyFormat.
	if err := aead.ValidateAESKeySize(format.AesCtrKeyFormat.KeySize); err != nil {
		return fmt.Errorf("aes_ctr_hmac_aead_key_manager: %s", err)
	}
	if format.AesCtrKeyFormat.Params.IvSize < aead.AESCTRMinIVSize || format.AesCtrKeyFormat.Params.IvSize > 16 {
		return errors.New("aes_ctr_hmac_aead_key_manager: invalid AesCtrHmacAeadKeyFormat: IV size out of range")
	}

	// Validate HmacKeyFormat.
	hmacKeyFormat := format.HmacKeyFormat
	if hmacKeyFormat.KeySize < minHMACKeySizeInBytes {
		return errors.New("aes_ctr_hmac_aead_key_manager: HMAC KeySize is too small")
	}
	if hmacKeyFormat.Params.TagSize < minTagSizeInBytes {
		return fmt.Errorf("aes_ctr_hmac_aead_key_manager: invalid HmacParams: TagSize %d is too small", hmacKeyFormat.Params.TagSize)
	}

	maxTagSize := map[commonpb.HashType]uint32{
		commonpb.HashType_SHA1:   20,
		commonpb.HashType_SHA256: 32,
		commonpb.HashType_SHA512: 64}

	tagSize, ok := maxTagSize[hmacKeyFormat.Params.Hash]
	if !ok {
		return fmt.Errorf("aes_ctr_hmac_aead_key_manager: invalid HmacParams: HashType %q not supported",
			hmacKeyFormat.Params.Hash)
	}
	if hmacKeyFormat.Params.TagSize > tagSize {
		return fmt.Errorf("aes_ctr_hmac_aead_key_manager: invalid HmacParams: tagSize %d is too big for HashType %q",
			hmacKeyFormat.Params.TagSize, hmacKeyFormat.Params.Hash)
	}

	return nil
}
