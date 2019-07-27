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

package hybrid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/tsingson/tink/golang/core/registry"
	"github.com/tsingson/tink/golang/keyset"
	subtle "github.com/tsingson/tink/golang/subtle/hybrid"
	eahpb "github.com/tsingson/tink/proto/ecies_aead_hkdf_go_proto"
	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

const (
	eciesAEADHKDFPrivateKeyKeyVersion = 0
	eciesAEADHKDFPrivateKeyTypeURL    = "type.googleapis.com/google.crypto.tink.EciesAeadHkdfPrivateKey"
)

// common errors
var errInvalidECIESAEADHKDFPrivateKeyKey = fmt.Errorf("ecies_aead_hkdf_private_key_manager: invalid key")
var errInvalidECIESAEADHKDFPrivateKeyKeyFormat = fmt.Errorf("ecies_aead_hkdf_private_key_manager: invalid key format")

// eciesAEADHKDFPrivateKeyKeyManager is an implementation of PrivateKeyManager interface.
// It generates new ECIESAEADHKDFPrivateKeyKey keys and produces new instances of ECIESAEADHKDFPrivateKey subtle.
type eciesAEADHKDFPrivateKeyKeyManager struct{}

// Assert that eciesAEADHKDFPrivateKeyKeyManager implements the PrivateKeyManager interface.
var _ registry.PrivateKeyManager = (*eciesAEADHKDFPrivateKeyKeyManager)(nil)

// newECIESAEADHKDFPrivateKeyKeyManager creates a new aesGcmKeyManager.
func newECIESAEADHKDFPrivateKeyKeyManager() *eciesAEADHKDFPrivateKeyKeyManager {
	return new(eciesAEADHKDFPrivateKeyKeyManager)
}

// Primitive creates an ECIESAEADHKDFPrivateKey subtle for the given serialized ECIESAEADHKDFPrivateKey proto.
func (km *eciesAEADHKDFPrivateKeyKeyManager) Primitive(serializedKey []byte) (interface{}, error) {
	if len(serializedKey) == 0 {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKey
	}
	key := new(eahpb.EciesAeadHkdfPrivateKey)
	if err := proto.Unmarshal(serializedKey, key); err != nil {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKey
	}
	if err := km.validateKey(key); err != nil {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKey
	}
	curve, err := subtle.GetCurve(key.PublicKey.Params.KemParams.CurveType.String())
	if err != nil {
		return nil, err
	}
	pvt := subtle.GetECPrivateKey(curve, key.KeyValue)
	rDem, err := newRegisterECIESAEADHKDFDemHelper(key.PublicKey.Params.DemParams.AeadDem)
	if err != nil {
		return nil, err
	}
	salt := key.PublicKey.Params.KemParams.HkdfSalt
	hash := key.PublicKey.Params.KemParams.HkdfHashType.String()
	ptFormat := key.PublicKey.Params.EcPointFormat.String()
	return subtle.NewECIESAEADHKDFHybridDecrypt(pvt, salt, hash, ptFormat, rDem)
}

// NewKey creates a new key according to specification the given serialized ECIESAEADHKDFPrivateKeyKeyFormat.
func (km *eciesAEADHKDFPrivateKeyKeyManager) NewKey(serializedKeyFormat []byte) (proto.Message, error) {
	if len(serializedKeyFormat) == 0 {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKeyFormat
	}
	keyFormat := new(eahpb.EciesAeadHkdfKeyFormat)
	if err := proto.Unmarshal(serializedKeyFormat, keyFormat); err != nil {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKeyFormat
	}
	if err := km.validateKeyFormat(keyFormat); err != nil {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKeyFormat
	}
	curve, err := subtle.GetCurve(keyFormat.Params.KemParams.CurveType.String())
	pvt, err := subtle.GenerateECDHKeyPair(curve)
	if err != nil {
		return nil, err
	}

	return &eahpb.EciesAeadHkdfPrivateKey{
		Version:  eciesAEADHKDFPrivateKeyKeyVersion,
		KeyValue: pvt.D.Bytes(),
		PublicKey: &eahpb.EciesAeadHkdfPublicKey{
			Version: eciesAEADHKDFPrivateKeyKeyVersion,
			Params:  keyFormat.Params,
			X:       pvt.PublicKey.Point.X.Bytes(),
			Y:       pvt.PublicKey.Point.Y.Bytes(),
		},
	}, nil
}

// NewKeyData creates a new KeyData according to specification in the given serialized
// ECIESAEADHKDFPrivateKeyKeyFormat.
// It should be used solely by the key management API.
func (km *eciesAEADHKDFPrivateKeyKeyManager) NewKeyData(serializedKeyFormat []byte) (*tinkpb.KeyData, error) {
	key, err := km.NewKey(serializedKeyFormat)
	if err != nil {
		return nil, err
	}
	serializedKey, err := proto.Marshal(key)
	if err != nil {
		return nil, err
	}
	return &tinkpb.KeyData{
		TypeUrl:         eciesAEADHKDFPrivateKeyTypeURL,
		Value:           serializedKey,
		KeyMaterialType: tinkpb.KeyData_ASYMMETRIC_PRIVATE,
	}, nil
}

func (km *eciesAEADHKDFPrivateKeyKeyManager) PublicKeyData(serializedPrivKey []byte) (*tinkpb.KeyData, error) {
	privKey := new(eahpb.EciesAeadHkdfPrivateKey)
	if err := proto.Unmarshal(serializedPrivKey, privKey); err != nil {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKey
	}
	serializedPubKey, err := proto.Marshal(privKey.PublicKey)
	if err != nil {
		return nil, errInvalidECIESAEADHKDFPrivateKeyKey
	}
	return &tinkpb.KeyData{
		TypeUrl:         eciesAEADHKDFPublicKeyTypeURL,
		Value:           serializedPubKey,
		KeyMaterialType: tinkpb.KeyData_ASYMMETRIC_PUBLIC,
	}, nil
}

// DoesSupport indicates if this key manager supports the given key type.
func (km *eciesAEADHKDFPrivateKeyKeyManager) DoesSupport(typeURL string) bool {
	return typeURL == eciesAEADHKDFPrivateKeyTypeURL
}

// TypeURL returns the key type of keys managed by this key manager.
func (km *eciesAEADHKDFPrivateKeyKeyManager) TypeURL() string {
	return eciesAEADHKDFPrivateKeyTypeURL
}

// validateKey validates the given ECDSAPrivateKey.
func (km *eciesAEADHKDFPrivateKeyKeyManager) validateKey(key *eahpb.EciesAeadHkdfPrivateKey) error {
	if err := keyset.ValidateKeyVersion(key.Version, eciesAEADHKDFPrivateKeyKeyVersion); err != nil {
		return fmt.Errorf("ecies_aead_hkdf_private_key_manager: invalid key: %s", err)
	}
	return checkECIESAEADHKDFParams(key.PublicKey.Params)
}

// validateKeyFormat validates the given ECDSAKeyFormat.
func (km *eciesAEADHKDFPrivateKeyKeyManager) validateKeyFormat(format *eahpb.EciesAeadHkdfKeyFormat) error {
	return checkECIESAEADHKDFParams(format.Params)
}

func checkECIESAEADHKDFParams(params *eahpb.EciesAeadHkdfParams) error {
	_, err := subtle.GetCurve(params.KemParams.CurveType.String())
	if err != nil {
		return err
	}
	if strings.Compare(params.KemParams.HkdfHashType.String(), "HashType_UNKNOWN_HASH") == 0 {
		return errors.New("hash unsupported for HMAC")
	}

	if strings.Compare(params.EcPointFormat.String(), "EcPointFormat_UNKNOWN_FORMAT") == 0 {
		return errors.New("unknown EC point format")
	}
	km, err := registry.GetKeyManager(params.DemParams.AeadDem.TypeUrl)
	if err != nil {
		return err
	}
	_, err = km.NewKeyData(params.DemParams.AeadDem.Value)
	if err != nil {
		return err
	}
	return nil
}
