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

package signature_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tsingson/tink/go/subtle"
	"github.com/tsingson/tink/go/subtle/random"
	subtleSignature "github.com/tsingson/tink/go/subtle/signature"
)

func TestSignVerify(t *testing.T) {
	data := random.GetRandomBytes(20)
	hash := "SHA256"
	curve := "NIST_P256"
	encodings := []string{"DER", "IEEE_P1363"}
	for _, encoding := range encodings {
		priv, _ := ecdsa.GenerateKey(subtle.GetCurve(curve), rand.Reader)
		// Use the private key and public key directly to create new instances
		signer, err := subtleSignature.NewECDSASignerFromPrivateKey(hash, encoding, priv)
		if err != nil {
			t.Errorf("unexpected error when creating ECDSASigner: %s", err)
		}
		verifier, err := subtleSignature.NewECDSAVerifierFromPublicKey(hash, encoding, &priv.PublicKey)
		if err != nil {
			t.Errorf("unexpected error when creating ECDSAVerifier: %s", err)
		}
		signature, err := signer.Sign(data)
		if err != nil {
			t.Errorf("unexpected error when signing: %s", err)
		}
		if err := verifier.Verify(signature, data); err != nil {
			t.Errorf("unexpected error when verifying: %s", err)
		}

		// Use byte slices to create new instances
		signer, err = subtleSignature.NewECDSASigner(hash, curve, encoding, priv.D.Bytes())
		if err != nil {
			t.Errorf("unexpected error when creating ECDSASigner: %s", err)
		}
		verifier, err = subtleSignature.NewECDSAVerifier(hash, curve, encoding, priv.X.Bytes(), priv.Y.Bytes())
		if err != nil {
			t.Errorf("unexpected error when creating ECDSAVerifier: %s", err)
		}
		signature, err = signer.Sign(data)
		if err != nil {
			t.Errorf("unexpected error when signing: %s", err)
		}
		if err = verifier.Verify(signature, data); err != nil {
			t.Errorf("unexpected error when verifying: %s", err)
		}
	}
}

type testData struct {
	Algorithm        string
	GeneratorVersion string
	NumberOfTests    uint32
	TestGroups       []*testGroup
}

type testGroup struct {
	Jwk    *jwk `json:"jwk,omitempty"`
	KeyDer string
	KeyPem string
	Sha    string
	Type   string
	Key    *testKey
	Tests  []*testcase
}

type testKey struct {
	Curve string
	Type  string
	Wx    string
	Wy    string
}

type testcase struct {
	Comment string
	Msg     string
	Result  string
	Sig     string
	TcID    uint32
}

type jwk struct {
	Jwk string
	Crv string
	Kid string
	Kty string
	X   string
	Y   string
}

func TestWycheproofVectors(t *testing.T) {
	vectors := []struct {
		Filename string
		Encoding string
	}{
		{"../../../third_party/wycheproof/testvectors/ecdsa_test.json", "DER"},
		{"../../../third_party/wycheproof/testvectors/ecdsa_secp256r1_sha256_p1363_test.json", "IEEE_P1363"},
		{"../../../third_party/wycheproof/testvectors/ecdsa_secp384r1_sha512_p1363_test.json", "IEEE_P1363"},
		{"../../../third_party/wycheproof/testvectors/ecdsa_secp521r1_sha512_p1363_test.json", "IEEE_P1363"},
	}

	for _, v := range vectors {
		wycheproofTest(t, v.Filename, v.Encoding)
	}
}

func wycheproofTest(t *testing.T, filename, encoding string) {
	fmt.Printf("now testing wycheproof file %q, encoding %q\n", filename, encoding)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot open file: %s, this is typically caused by an older version of Wycheproof.", err)
		return
	}
	parser := json.NewDecoder(f)
	content := new(testData)
	if err := parser.Decode(content); err != nil {
		t.Errorf("cannot decode content of file: %s", err)
	}
	for _, g := range content.TestGroups {
		hash := subtle.ConvertHashName(g.Sha)
		curve := subtle.ConvertCurveName(g.Key.Curve)
		if hash == "" || curve == "" {
			continue
		}
		x, err := subtle.NewBigIntFromHex(g.Key.Wx)
		if err != nil {
			t.Errorf("cannot decode wx: %s", err)
		}
		y, err := subtle.NewBigIntFromHex(g.Key.Wy)
		if err != nil {
			t.Errorf("cannot decode wy: %s", err)
		}
		verifier, err := subtleSignature.NewECDSAVerifier(hash, curve, encoding, x.Bytes(), y.Bytes())
		if err != nil {
			continue
		}
		for _, tc := range g.Tests {
			message, err := hex.DecodeString(tc.Msg)
			if err != nil {
				t.Errorf("cannot decode message in test case %d: %s", tc.TcID, err)
			}
			sig, err := hex.DecodeString(tc.Sig)
			if err != nil {
				t.Errorf("cannot decode signature in test case %d: %s", tc.TcID, err)
			}
			err = verifier.Verify(sig, message)
			if (tc.Result == "valid" && err != nil) ||
				(tc.Result == "invalid" && err == nil) {
				t.Errorf("failed in test case %d with error %q ", tc.TcID, err)
			}
		}
	}
}
