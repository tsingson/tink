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

// Package daead provides implementations of the DeterministicAEAD primitive.
// Unlike AEAD, implementations of this interface are not semantically secure, because
// encrypting the same plaintex always yields the same ciphertext.
// Example:
//
// package main
//
// import (
//     "fmt"
//
//     "github.com/tsingson/tink/go/daead"
//     "github.com/tsingson/tink/go/keyset"
// )
//
// func main() {
//
//     kh, err := keyset.NewHandle(daead.AESSIVKeyTemplate())
//     if err != nil {
//         // handle the error
//     }
//
//     d := daead.New(kh)
//
//     ct1 , err := d.EncryptDeterministically([]byte("this data needs to be encrypted"), []byte("additional data"))
//     if err != nil {
//         // handle error
//     }
//
//     pt , err := d.DecryptDeterministically(ct, []byte("additional data"))
//     if err != nil {
//         // handle error
//     }
//
//     ct2 , err := d.EncryptDeterministically([]byte("this data needs to be encrypted"), []byte("additional data"))
//     if err != nil {
//         // handle error
//     }
//
//     // ct1 will be equal to ct2
//
// }
package daead

import (
	"fmt"

	"github.com/tsingson/tink/go/core/registry"
)

func init() {
	if err := registry.RegisterKeyManager(newAESSIVKeyManager()); err != nil {
		panic(fmt.Sprintf("daead.init() failed: %v", err))
	}
}
