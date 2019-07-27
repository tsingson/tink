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

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/tsingson/tink/golang/subtle/random"
	"github.com/tsingson/tink/golang/tink"
)

// XChaCha20Poly1305 is an implementation of AEAD interface.
type XChaCha20Poly1305 struct {
	Key []byte
}

// Assert that XChaCha20Poly1305 implements the AEAD interface.
var _ tink.AEAD = (*XChaCha20Poly1305)(nil)

// NewXChaCha20Poly1305 returns an XChaCha20Poly1305 instance.
// The key argument should be a 32-bytes key.
func NewXChaCha20Poly1305(key []byte) (*XChaCha20Poly1305, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("xchacha20poly1305: bad key length")
	}

	return &XChaCha20Poly1305{Key: key}, nil
}

// Encrypt encrypts {@code pt} with {@code aad} as additional
// authenticated data. The resulting ciphertext consists of two parts:
// (1) the nonce used for encryption and (2) the actual ciphertext.
func (x *XChaCha20Poly1305) Encrypt(pt []byte, aad []byte) ([]byte, error) {
	c, err := chacha20poly1305.NewX(x.Key)
	if err != nil {
		return nil, err
	}

	n := x.newNonce()
	ct := c.Seal(nil, n, pt, aad)
	var ret []byte
	ret = append(ret, n...)
	ret = append(ret, ct...)
	return ret, nil
}

// Decrypt decrypts {@code ct} with {@code aad} as the additionalauthenticated data.
func (x *XChaCha20Poly1305) Decrypt(ct []byte, aad []byte) ([]byte, error) {
	c, err := chacha20poly1305.NewX(x.Key)
	if err != nil {
		return nil, err
	}

	n := ct[:chacha20poly1305.NonceSizeX]
	pt, err := c.Open(nil, n, ct[chacha20poly1305.NonceSizeX:], aad)
	if err != nil {
		return nil, fmt.Errorf("XChaCha20Poly1305.Decrypt: %s", err)
	}
	return pt, nil
}

// newNonce creates a new nonce for encryption.
func (x *XChaCha20Poly1305) newNonce() []byte {
	return random.GetRandomBytes(chacha20poly1305.NonceSizeX)
}
