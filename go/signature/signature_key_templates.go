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
	"github.com/golang/protobuf/proto"

	commonpb "github.com/tsingson/tink/proto/common_go_proto"
	ecdsapb "github.com/tsingson/tink/proto/ecdsa_go_proto"
	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

// This file contains pre-generated KeyTemplates for Signer and Verifier.
// One can use these templates to generate new Keysets.

// ECDSAP256KeyTemplate is a KeyTemplate that generates a new ECDSA private key with the following parameters:
//   - Hash function: SHA256
//   - Curve: NIST P-256
//   - Signature encoding: DER
func ECDSAP256KeyTemplate() *tinkpb.KeyTemplate {
	return createECDSAKeyTemplate(commonpb.HashType_SHA256,
		commonpb.EllipticCurveType_NIST_P256,
		ecdsapb.EcdsaSignatureEncoding_DER)
}

// ECDSAP384KeyTemplate is a KeyTemplate that generates a new ECDSA private key with the following parameters:
//   - Hash function: SHA512
//   - Curve: NIST P-384
//   - Signature encoding: DER
func ECDSAP384KeyTemplate() *tinkpb.KeyTemplate {
	return createECDSAKeyTemplate(commonpb.HashType_SHA512,
		commonpb.EllipticCurveType_NIST_P384,
		ecdsapb.EcdsaSignatureEncoding_DER)
}

// ECDSAP521KeyTemplate is a KeyTemplate that generates a new ECDSA private key with the following parameters:
//   - Hash function: SHA512
//   - Curve: NIST P-521
//   - Signature encoding: DER
func ECDSAP521KeyTemplate() *tinkpb.KeyTemplate {
	return createECDSAKeyTemplate(commonpb.HashType_SHA512,
		commonpb.EllipticCurveType_NIST_P521,
		ecdsapb.EcdsaSignatureEncoding_DER)
}

// createECDSAKeyTemplate creates a KeyTemplate containing a EcdasKeyFormat
// with the given parameters.
func createECDSAKeyTemplate(hashType commonpb.HashType, curve commonpb.EllipticCurveType, encoding ecdsapb.EcdsaSignatureEncoding) *tinkpb.KeyTemplate {
	params := &ecdsapb.EcdsaParams{
		HashType: hashType,
		Curve:    curve,
		Encoding: encoding,
	}
	format := &ecdsapb.EcdsaKeyFormat{Params: params}
	serializedFormat, _ := proto.Marshal(format)
	return &tinkpb.KeyTemplate{
		TypeUrl: ecdsaSignerTypeURL,
		Value:   serializedFormat,
	}
}

// ED25519KeyTemplate is a KeyTemplate that generates a new ED25519 private key.
func ED25519KeyTemplate() *tinkpb.KeyTemplate {
	return &tinkpb.KeyTemplate{
		TypeUrl: ed25519SignerTypeURL,
	}
}
