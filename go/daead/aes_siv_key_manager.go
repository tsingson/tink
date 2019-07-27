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

package daead

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/tsingson/tink/go/core/registry"
	"github.com/tsingson/tink/go/keyset"
	"github.com/tsingson/tink/go/subtle/daead"
	"github.com/tsingson/tink/go/subtle/random"

	aspb "github.com/tsingson/tink/proto/aes_siv_go_proto"
	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

const (
	aesSIVKeyVersion = 0
	aesSIVTypeURL    = "type.googleapis.com/google.crypto.tink.AesSivKey"
)

// aesSIVKeyManager is an implementation of KeyManager interface.
// It generates new AesSivKey keys and produces new instances of AESSIV subtle.
type aesSIVKeyManager struct{}

// Assert that aesSIVKeyManager implements the KeyManager interface.
var _ registry.KeyManager = (*aesSIVKeyManager)(nil)

// newAESSIVKeyManager creates a new aesSIVKeyManager.
func newAESSIVKeyManager() *aesSIVKeyManager {
	return new(aesSIVKeyManager)
}

// Primitive creates an AESSIV subtle for the given serialized AesSivKey proto.
func (km *aesSIVKeyManager) Primitive(serializedKey []byte) (interface{}, error) {
	if len(serializedKey) == 0 {
		return nil, fmt.Errorf("aes_siv_key_manager: invalid key")
	}

	key := new(aspb.AesSivKey)
	if err := proto.Unmarshal(serializedKey, key); err != nil {
		return nil, err
	}
	if err := km.validateKey(key); err != nil {
		return nil, err
	}
	ret, err := daead.NewAESSIV(key.KeyValue)
	if err != nil {
		return nil, fmt.Errorf("aes_siv_key_manager: cannot create new primitive: %s", err)
	}
	return ret, nil
}

// NewKey creates a new key, ignoring the specification in the given serialized key format
// because the key size and other params are fixed.
func (km *aesSIVKeyManager) NewKey(serializedKeyFormat []byte) (proto.Message, error) {
	return km.newAesSivKey(), nil
}

// NewKeyData creates a new KeyData, ignoring the specification in the given serialized key format
// because the key size and other params are fixed.
// It should be used solely by the key management API.
func (km *aesSIVKeyManager) NewKeyData(serializedKeyFormat []byte) (*tinkpb.KeyData, error) {
	key := km.newAesSivKey()
	serializedKey, err := proto.Marshal(key)
	if err != nil {
		return nil, err
	}
	return &tinkpb.KeyData{
		TypeUrl:         aesSIVTypeURL,
		Value:           serializedKey,
		KeyMaterialType: tinkpb.KeyData_SYMMETRIC,
	}, nil
}

// DoesSupport indicates if this key manager supports the given key type.
func (km *aesSIVKeyManager) DoesSupport(typeURL string) bool {
	return typeURL == aesSIVTypeURL
}

// TypeURL returns the key type of keys managed by this key manager.
func (km *aesSIVKeyManager) TypeURL() string {
	return aesSIVTypeURL
}

// validateKey validates the given AesSivKey.
func (km *aesSIVKeyManager) validateKey(key *aspb.AesSivKey) error {
	err := keyset.ValidateKeyVersion(key.Version, aesSIVKeyVersion)
	if err != nil {
		return fmt.Errorf("aes_siv_key_manager: %s", err)
	}
	keySize := uint32(len(key.KeyValue))
	if keySize != daead.AESSIVKeySize {
		return fmt.Errorf("aes_siv_key_manager: keySize != %d", daead.AESSIVKeySize)
	}
	return nil
}

// newAesSivKey creates a new AesSivKey.
func (km *aesSIVKeyManager) newAesSivKey() *aspb.AesSivKey {
	keyValue := random.GetRandomBytes(daead.AESSIVKeySize)
	return &aspb.AesSivKey{
		Version:  aesSIVKeyVersion,
		KeyValue: keyValue,
	}
}
