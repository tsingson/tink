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

package keyset_test

import (
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/tsingson/tink/go/keyset"
	"github.com/tsingson/tink/go/mac"
	"github.com/tsingson/tink/go/subtle/aead"
	"github.com/tsingson/tink/go/testkeyset"
	"github.com/tsingson/tink/go/testutil"

	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

func TestNewHandle(t *testing.T) {
	kt := mac.HMACSHA256Tag128KeyTemplate()
	kh, err := keyset.NewHandle(kt)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	ks := testkeyset.KeysetMaterial(kh)
	if len(ks.Key) != 1 {
		t.Errorf("incorrect number of keys in the keyset: %d", len(ks.Key))
	}
	key := ks.Key[0]
	if ks.PrimaryKeyId != key.KeyId {
		t.Errorf("incorrect primary key id, expect %d, got %d", key.KeyId, ks.PrimaryKeyId)
	}
	if key.KeyData.TypeUrl != kt.TypeUrl {
		t.Errorf("incorrect type url, expect %s, got %s", kt.TypeUrl, key.KeyData.TypeUrl)
	}
	if _, err = mac.New(kh); err != nil {
		t.Errorf("cannot get primitive from generated keyset handle: %s", err)
	}
}

func TestNewHandleWithInvalidInput(t *testing.T) {
	// template unregistered TypeUrl
	template := mac.HMACSHA256Tag128KeyTemplate()
	template.TypeUrl = "some unknown TypeUrl"
	if _, err := keyset.NewHandle(template); err == nil {
		t.Errorf("expect an error when TypeUrl is not registered")
	}
	// nil
	if _, err := keyset.NewHandle(nil); err == nil {
		t.Errorf("expect an error when template is nil")
	}
}

func TestRead(t *testing.T) {
	masterKey, err := aead.NewAESGCM([]byte(strings.Repeat("A", 32)))
	if err != nil {
		t.Errorf("aead.NewAESGCM(): %v", err)
	}

	// Create a keyset
	keyData := testutil.NewKeyData("some type url", []byte{0}, tinkpb.KeyData_SYMMETRIC)
	key := testutil.NewKey(keyData, tinkpb.KeyStatusType_ENABLED, 1, tinkpb.OutputPrefixType_TINK)
	ks := testutil.NewKeyset(1, []*tinkpb.Keyset_Key{key})
	h, _ := testkeyset.NewHandle(ks)

	memKeyset := &keyset.MemReaderWriter{}
	if err := h.Write(memKeyset, masterKey); err != nil {
		t.Fatalf("handle.Write(): %v", err)
	}
	h2, err := keyset.Read(memKeyset, masterKey)
	if err != nil {
		t.Fatalf("keyset.Read(): %v", err)
	}
	if !proto.Equal(testkeyset.KeysetMaterial(h), testkeyset.KeysetMaterial(h2)) {
		t.Fatalf("Decrypt failed: got %v, want %v", h2, h)
	}
}

func TestReadWithNoSecrets(t *testing.T) {
	// Create a keyset containing public key material
	keyData := testutil.NewKeyData("some type url", []byte{0}, tinkpb.KeyData_ASYMMETRIC_PUBLIC)
	key := testutil.NewKey(keyData, tinkpb.KeyStatusType_ENABLED, 1, tinkpb.OutputPrefixType_TINK)
	ks := testutil.NewKeyset(1, []*tinkpb.Keyset_Key{key})
	h, _ := testkeyset.NewHandle(ks)

	memKeyset := &keyset.MemReaderWriter{}
	if err := h.WriteWithNoSecrets(memKeyset); err != nil {
		t.Fatalf("handle.WriteWithNoSecrets(): %v", err)
	}
	h2, err := keyset.ReadWithNoSecrets(memKeyset)
	if err != nil {
		t.Fatalf("keyset.ReadWithNoSecrets(): %v", err)
	}
	if !proto.Equal(testkeyset.KeysetMaterial(h), testkeyset.KeysetMaterial(h2)) {
		t.Fatalf("Decrypt failed: got %v, want %v", h2, h)
	}
}

func TestWithNoSecretsFunctionsFailWhenHandlingSecretKeyMaterial(t *testing.T) {
	// Create a keyset containing secret key material (symmetric)
	keyData := testutil.NewKeyData("some type url", []byte{0}, tinkpb.KeyData_SYMMETRIC)
	key := testutil.NewKey(keyData, tinkpb.KeyStatusType_ENABLED, 1, tinkpb.OutputPrefixType_TINK)
	ks := testutil.NewKeyset(1, []*tinkpb.Keyset_Key{key})
	h, _ := testkeyset.NewHandle(ks)

	if err := h.WriteWithNoSecrets(&keyset.MemReaderWriter{}); err == nil {
		t.Error("handle.WriteWithNoSecrets() should fail when exporting secret key material")
	}

	if _, err := keyset.ReadWithNoSecrets(&keyset.MemReaderWriter{Keyset: testkeyset.KeysetMaterial(h)}); err == nil {
		t.Error("keyset.ReadWithNoSecrets should fail when importing secret key material")
	}
}

func TestWithNoSecretsFunctionsFailWhenUnknownKeyMaterial(t *testing.T) {
	// Create a keyset containing secret key material (symmetric)
	keyData := testutil.NewKeyData("some type url", []byte{0}, tinkpb.KeyData_UNKNOWN_KEYMATERIAL)
	key := testutil.NewKey(keyData, tinkpb.KeyStatusType_ENABLED, 1, tinkpb.OutputPrefixType_TINK)
	ks := testutil.NewKeyset(1, []*tinkpb.Keyset_Key{key})
	h, _ := testkeyset.NewHandle(ks)

	if err := h.WriteWithNoSecrets(&keyset.MemReaderWriter{}); err == nil {
		t.Error("handle.WriteWithNoSecrets() should fail when exporting secret key material")
	}

	if _, err := keyset.ReadWithNoSecrets(&keyset.MemReaderWriter{Keyset: testkeyset.KeysetMaterial(h)}); err == nil {
		t.Error("keyset.ReadWithNoSecrets should fail when importing secret key material")
	}
}

func TestWithNoSecretsFunctionsFailWithAsymmetricPrivateKeyMaterial(t *testing.T) {
	// Create a keyset containing secret key material (symmetric)
	keyData := testutil.NewKeyData("some type url", []byte{0}, tinkpb.KeyData_ASYMMETRIC_PRIVATE)
	key := testutil.NewKey(keyData, tinkpb.KeyStatusType_ENABLED, 1, tinkpb.OutputPrefixType_TINK)
	ks := testutil.NewKeyset(1, []*tinkpb.Keyset_Key{key})
	h, _ := testkeyset.NewHandle(ks)

	if err := h.WriteWithNoSecrets(&keyset.MemReaderWriter{}); err == nil {
		t.Error("handle.WriteWithNoSecrets() should fail when exporting secret key material")
	}

	if _, err := keyset.ReadWithNoSecrets(&keyset.MemReaderWriter{Keyset: testkeyset.KeysetMaterial(h)}); err == nil {
		t.Error("keyset.ReadWithNoSecrets should fail when importing secret key material")
	}
}
