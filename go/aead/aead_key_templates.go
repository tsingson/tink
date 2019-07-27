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
	"github.com/golang/protobuf/proto"

	ctrpb "github.com/tsingson/tink/proto/aes_ctr_go_proto"
	ctrhmacpb "github.com/tsingson/tink/proto/aes_ctr_hmac_aead_go_proto"
	gcmpb "github.com/tsingson/tink/proto/aes_gcm_go_proto"
	commonpb "github.com/tsingson/tink/proto/common_go_proto"
	hmacpb "github.com/tsingson/tink/proto/hmac_go_proto"
	kmsenvpb "github.com/tsingson/tink/proto/kms_envelope_go_proto"
	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

// This file contains pre-generated KeyTemplates for AEAD keys. One can use these templates
// to generate new Keysets.

// AES128GCMKeyTemplate is a KeyTemplate that generates an AES-GCM key with the following parameters:
//   - Key size: 16 bytes
func AES128GCMKeyTemplate() *tinkpb.KeyTemplate {
	return createAESGCMKeyTemplate(16)
}

// AES256GCMKeyTemplate is a KeyTemplate that generates an AES-GCM key with the following parameters:
//   - Key size: 32 bytes
func AES256GCMKeyTemplate() *tinkpb.KeyTemplate {
	return createAESGCMKeyTemplate(32)
}

// AES128CTRHMACSHA256KeyTemplate is a KeyTemplate that generates an AES-CTR-HMAC-AEAD key with the following parameters:
//  - AES key size: 16 bytes
//  - AES CTR IV size: 16 bytes
//  - HMAC key size: 32 bytes
//  - HMAC tag size: 16 bytes
//  - HMAC hash function: SHA256
func AES128CTRHMACSHA256KeyTemplate() *tinkpb.KeyTemplate {
	return createAESCTRHMACAEADKeyTemplate(16, 16, 32, 16, commonpb.HashType_SHA256)
}

// AES256CTRHMACSHA256KeyTemplate is a KeyTemplate that generates an AES-CTR-HMAC-AEAD key with the following parameters:
//  - AES key size: 32 bytes
//  - AES CTR IV size: 16 bytes
//  - HMAC key size: 32 bytes
//  - HMAC tag size: 32 bytes
//  - HMAC hash function: SHA256
func AES256CTRHMACSHA256KeyTemplate() *tinkpb.KeyTemplate {
	return createAESCTRHMACAEADKeyTemplate(32, 16, 32, 32, commonpb.HashType_SHA256)
}

// ChaCha20Poly1305KeyTemplate is a KeyTemplate that generates a CHACHA20_POLY1305 key.
func ChaCha20Poly1305KeyTemplate() *tinkpb.KeyTemplate {
	return &tinkpb.KeyTemplate{
		// Don't set value because KeyFormat is not required.
		TypeUrl:          chaCha20Poly1305TypeURL,
		OutputPrefixType: tinkpb.OutputPrefixType_TINK,
	}
}

// XChaCha20Poly1305KeyTemplate is a KeyTemplate that generates a XCHACHA20_POLY1305 key.
func XChaCha20Poly1305KeyTemplate() *tinkpb.KeyTemplate {
	return &tinkpb.KeyTemplate{
		// Don't set value because KeyFormat is not required.
		TypeUrl:          xChaCha20Poly1305TypeURL,
		OutputPrefixType: tinkpb.OutputPrefixType_TINK,
	}
}

// KMSEnvelopeAEADKeyTemplate is a KeyTemplate that generates a KMSEnvelopeAEAD key for a given KEK in remote KMS
func KMSEnvelopeAEADKeyTemplate(uri string, dekT *tinkpb.KeyTemplate) *tinkpb.KeyTemplate {
	f := &kmsenvpb.KmsEnvelopeAeadKeyFormat{
		KekUri:      uri,
		DekTemplate: dekT,
	}
	serializedFormat, _ := proto.Marshal(f)
	return &tinkpb.KeyTemplate{
		Value:            serializedFormat,
		TypeUrl:          kmsEnvelopeAEADTypeURL,
		OutputPrefixType: tinkpb.OutputPrefixType_TINK,
	}
}

// createAESGCMKeyTemplate creates a new AES-GCM key template with the given key
// size in bytes.
func createAESGCMKeyTemplate(keySize uint32) *tinkpb.KeyTemplate {
	format := &gcmpb.AesGcmKeyFormat{
		KeySize: keySize,
	}
	serializedFormat, _ := proto.Marshal(format)
	return &tinkpb.KeyTemplate{
		TypeUrl: aesGCMTypeURL,
		Value:   serializedFormat,
	}
}

func createAESCTRHMACAEADKeyTemplate(aesKeySize, ivSize, hmacKeySize, tagSize uint32, hash commonpb.HashType) *tinkpb.KeyTemplate {
	format := &ctrhmacpb.AesCtrHmacAeadKeyFormat{
		AesCtrKeyFormat: &ctrpb.AesCtrKeyFormat{
			Params:  &ctrpb.AesCtrParams{IvSize: ivSize},
			KeySize: aesKeySize,
		},
		HmacKeyFormat: &hmacpb.HmacKeyFormat{
			Params:  &hmacpb.HmacParams{Hash: hash, TagSize: tagSize},
			KeySize: hmacKeySize,
		},
	}
	serializedFormat, _ := proto.Marshal(format)
	return &tinkpb.KeyTemplate{
		Value:            serializedFormat,
		TypeUrl:          aesCTRHMACAEADTypeURL,
		OutputPrefixType: tinkpb.OutputPrefixType_TINK,
	}
}
