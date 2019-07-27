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
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/tsingson/tink/golang/core/registry"
	"github.com/tsingson/tink/golang/tink"
	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

const (
	lenDEK = 4
)

// KMSEnvelopeAEAD represents an instance of Envelope AEAD.
type KMSEnvelopeAEAD struct {
	dekTemplate *tinkpb.KeyTemplate
	remote      tink.AEAD
}

var _ tink.AEAD = (*KMSEnvelopeAEAD)(nil)

// NewKMSEnvelopeAEAD creates an new instance of KMSEnvelopeAEAD
func NewKMSEnvelopeAEAD(kt tinkpb.KeyTemplate, remote tink.AEAD) *KMSEnvelopeAEAD {
	return &KMSEnvelopeAEAD{
		remote:      remote,
		dekTemplate: &kt,
	}
}

// Encrypt implements the tink.AEAD interface for encryption.
func (a *KMSEnvelopeAEAD) Encrypt(pt, aad []byte) ([]byte, error) {
	dekM, err := registry.NewKey(a.dekTemplate)
	if err != nil {
		return nil, err
	}
	dek, err := proto.Marshal(dekM)
	if err != nil {
		return nil, err
	}
	encryptedDEK, err := a.remote.Encrypt(dek, []byte{})
	if err != nil {
		return nil, err
	}
	p, err := registry.Primitive(a.dekTemplate.TypeUrl, dek)
	if err != nil {
		return nil, err
	}
	primitive, ok := p.(tink.AEAD)
	if !ok {
		return nil, errors.New("kms_envelope_aead: failed to convert AEAD primitive")
	}

	payload, err := primitive.Encrypt(pt, aad)
	if err != nil {
		return nil, err
	}
	return buildCipherText(encryptedDEK, payload)

}

// Decrypt implements the tink.AEAD interface for decryption.
func (a *KMSEnvelopeAEAD) Decrypt(ct, aad []byte) ([]byte, error) {
	b := bytes.NewBuffer(ct)
	bLen := b.Len()

	ed := int(binary.BigEndian.Uint32(b.Next(lenDEK)))
	if ed <= 0 || ed > len(ct)-lenDEK {
		return nil, errors.New("kms_envelope_aead: invalid ciphertext")
	}

	encryptedDEK := make([]byte, ed)
	n, err := b.Read(encryptedDEK)
	if err != nil || n != ed {
		return nil, errors.New("kms_envelope_aead: invalid ciphertext")
	}

	pl := bLen - lenDEK - ed
	payload := make([]byte, pl)
	n, err = b.Read(payload)
	if err != nil || n != pl {
		return nil, errors.New("kms_envelope_aead: invalid ciphertext")
	}

	dek, err := a.remote.Decrypt(encryptedDEK, []byte{})
	if err != nil {
		return nil, err
	}
	p, err := registry.Primitive(a.dekTemplate.TypeUrl, dek)
	if err != nil {
		return nil, fmt.Errorf("kms_envelope_aead: %s", err)
	}
	primitive, ok := p.(tink.AEAD)
	if !ok {
		return nil, errors.New("kms_envelope_aead: failed to convert AEAD primitive")
	}
	return primitive.Decrypt(payload, aad)
}

// buildCipherText builds the cipher text by appending the length DEK, encrypted DEK
// and the encrypted payload.
func buildCipherText(encryptedDEK, payload []byte) ([]byte, error) {
	var b bytes.Buffer

	// Write the length of the encrypted DEK.
	lenDEKbuf := make([]byte, lenDEK)
	binary.BigEndian.PutUint32(lenDEKbuf, uint32(len(encryptedDEK)))
	_, err := b.Write(lenDEKbuf)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(encryptedDEK)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(payload)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
