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

	"github.com/tsingson/tink/go/core/cryptofmt"
	"github.com/tsingson/tink/go/core/primitiveset"
	"github.com/tsingson/tink/go/core/registry"
	"github.com/tsingson/tink/go/keyset"
	"github.com/tsingson/tink/go/tink"
)

// New returns a DeterministicAEAD primitive from the given keyset handle.
func New(h *keyset.Handle) (tink.DeterministicAEAD, error) {
	return NewWithKeyManager(h, nil /*keyManager*/)
}

// NewWithKeyManager returns a DeterministicAEAD primitive from the given keyset handle and custom key manager.
func NewWithKeyManager(h *keyset.Handle, km registry.KeyManager) (tink.DeterministicAEAD, error) {
	ps, err := h.PrimitivesWithKeyManager(km)
	if err != nil {
		return nil, fmt.Errorf("daead_factory: cannot obtain primitive set: %s", err)
	}
	ret := new(primitiveSet)
	ret.ps = ps
	return tink.DeterministicAEAD(ret), nil
}

// primitiveSet is an DeterministicAEAD implementation that uses the underlying primitive set
// for deterministic encryption and decryption.
type primitiveSet struct {
	ps *primitiveset.PrimitiveSet
}

// Asserts that primitiveSet implements the DeterministicAEAD interface.
var _ tink.DeterministicAEAD = (*primitiveSet)(nil)

// EncryptDeterministically deterministically encrypts plaintext with additionalData as additional authenticated data.
// It returns the concatenation of the primary's identifier and the ciphertext.
func (d *primitiveSet) EncryptDeterministically(pt, aad []byte) ([]byte, error) {
	primary := d.ps.Primary
	p := (primary.Primitive).(tink.DeterministicAEAD)
	ct, err := p.EncryptDeterministically(pt, aad)
	if err != nil {
		return nil, err
	}

	var ret []byte
	ret = append(ret, primary.Prefix...)
	ret = append(ret, ct...)
	return ret, nil
}

// DecryptDeterministically deterministically decrypts ciphertext with additionalData as
// additional authenticated data. It returns the corresponding plaintext if the
// ciphertext is authenticated.
func (d *primitiveSet) DecryptDeterministically(ct, aad []byte) ([]byte, error) {
	// try non-raw keys
	prefixSize := cryptofmt.NonRawPrefixSize
	if len(ct) > prefixSize {
		prefix := ct[:prefixSize]
		ctNoPrefix := ct[prefixSize:]
		entries, err := d.ps.EntriesForPrefix(string(prefix))
		if err == nil {
			for i := 0; i < len(entries); i++ {
				p := (entries[i].Primitive).(tink.DeterministicAEAD)
				pt, err := p.DecryptDeterministically(ctNoPrefix, aad)
				if err == nil {
					return pt, nil
				}
			}
		}
	}
	// try raw keys
	entries, err := d.ps.RawEntries()
	if err == nil {
		for i := 0; i < len(entries); i++ {
			p := (entries[i].Primitive).(tink.DeterministicAEAD)
			pt, err := p.DecryptDeterministically(ct, aad)
			if err == nil {
				return pt, nil
			}
		}
	}
	// nothing worked
	return nil, fmt.Errorf("daead_factory: decryption failed")
}
