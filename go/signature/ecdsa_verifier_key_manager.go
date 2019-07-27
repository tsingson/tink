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

package signature

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/tsingson/tink/go/core/registry"
	"github.com/tsingson/tink/go/keyset"
	subtleSignature "github.com/tsingson/tink/go/subtle/signature"
	ecdsapb "github.com/tsingson/tink/proto/ecdsa_go_proto"
	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

const (
	ecdsaVerifierKeyVersion = 0
	ecdsaVerifierTypeURL    = "type.googleapis.com/google.crypto.tink.EcdsaPublicKey"
)

// common errors
var errInvalidECDSAVerifierKey = fmt.Errorf("ecdsa_verifier_key_manager: invalid key")
var errECDSAVerifierNotImplemented = fmt.Errorf("ecdsa_verifier_key_manager: not implemented")

// ecdsaVerifierKeyManager is an implementation of KeyManager interface.
// It doesn't support key generation.
type ecdsaVerifierKeyManager struct{}

// Assert that ecdsaVerifierKeyManager implements the KeyManager interface.
var _ registry.KeyManager = (*ecdsaVerifierKeyManager)(nil)

// newECDSAVerifierKeyManager creates a new ecdsaVerifierKeyManager.
func newECDSAVerifierKeyManager() *ecdsaVerifierKeyManager {
	return new(ecdsaVerifierKeyManager)
}

// Primitive creates an ECDSAVerifier subtle for the given serialized ECDSAPublicKey proto.
func (km *ecdsaVerifierKeyManager) Primitive(serializedKey []byte) (interface{}, error) {
	if len(serializedKey) == 0 {
		return nil, errInvalidECDSAVerifierKey
	}
	key := new(ecdsapb.EcdsaPublicKey)
	if err := proto.Unmarshal(serializedKey, key); err != nil {
		return nil, errInvalidECDSAVerifierKey
	}
	if err := km.validateKey(key); err != nil {
		return nil, fmt.Errorf("ecdsa_verifier_key_manager: %s", err)
	}
	hash, curve, encoding := getECDSAParamNames(key.Params)
	ret, err := subtleSignature.NewECDSAVerifier(hash, curve, encoding, key.X, key.Y)
	if err != nil {
		return nil, fmt.Errorf("ecdsa_verifier_key_manager: invalid key: %s", err)
	}
	return ret, nil
}

// NewKey is not implemented.
func (km *ecdsaVerifierKeyManager) NewKey(serializedKeyFormat []byte) (proto.Message, error) {
	return nil, errECDSAVerifierNotImplemented
}

// NewKeyData creates a new KeyData according to specification in  the given
// serialized ECDSAKeyFormat. It should be used solely by the key management API.
func (km *ecdsaVerifierKeyManager) NewKeyData(serializedKeyFormat []byte) (*tinkpb.KeyData, error) {
	return nil, errECDSAVerifierNotImplemented
}

// DoesSupport indicates if this key manager supports the given key type.
func (km *ecdsaVerifierKeyManager) DoesSupport(typeURL string) bool {
	return typeURL == ecdsaVerifierTypeURL
}

// TypeURL returns the key type of keys managed by this key manager.
func (km *ecdsaVerifierKeyManager) TypeURL() string {
	return ecdsaVerifierTypeURL
}

// validateKey validates the given ECDSAPublicKey.
func (km *ecdsaVerifierKeyManager) validateKey(key *ecdsapb.EcdsaPublicKey) error {
	if err := keyset.ValidateKeyVersion(key.Version, ecdsaVerifierKeyVersion); err != nil {
		return fmt.Errorf("ecdsa_verifier_key_manager: %s", err)
	}
	hash, curve, encoding := getECDSAParamNames(key.Params)
	return subtleSignature.ValidateECDSAParams(hash, curve, encoding)
}
