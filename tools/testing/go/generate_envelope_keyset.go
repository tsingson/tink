// Copyright 2017 Google LLC.
//
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

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"flag"
	// context is used to cancel outstanding requests
	// TEST_SRCDIR to read the roots.pem
	"github.com/tsingson/tink/go/aead"
	"github.com/tsingson/tink/go/core/registry"
	"github.com/tsingson/tink/go/insecurecleartextkeyset"
	"github.com/tsingson/tink/go/integration/awskms"
	"github.com/tsingson/tink/go/integration/gcpkms"
	"github.com/tsingson/tink/go/keyset"

	tinkpb "github.com/tsingson/tink/proto/tink_go_proto"
)

// lint placeholder header, please ignore
var (
	gcpURI      = "gcp-kms://projects/tink-test-infrastructure/locations/global/keyRings/unit-and-integration-testing/cryptoKeys/aead-key"
	gcpCredFile = os.Getenv("TEST_SRCDIR") + "/" + os.Getenv("TEST_WORKSPACE") + "/" + "testdata/credential.json"
	awsURI      = "aws-kms://arn:aws:kms:us-east-2:235739564943:key/3ee50705-5a82-4f5b-9753-05c4f473922f"
	awsCredFile = os.Getenv("TEST_SRCDIR") + "/" + os.Getenv("TEST_WORKSPACE") + "/" + "testdata/credentials_aws.csv"
)

func init() {
	certPath := os.Getenv("TEST_SRCDIR") + "/" + os.Getenv("TEST_WORKSPACE") + "/" + "roots.pem"
	flag.Set("cacerts", certPath)
	os.Setenv("SSL_CERT_FILE", certPath)
}

// lint placeholder footer, please ignore

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s keyset-file kms dek-template", os.Args[0])
	}
	f := os.Args[1]
	kms := os.Args[2]
	dek := os.Args[3]
	var dekT *tinkpb.KeyTemplate
	var kh *keyset.Handle
	var b bytes.Buffer
	switch strings.ToUpper(dek) {
	case "AES128_GCM":
		dekT = aead.AES128GCMKeyTemplate()
	case "AES128_CTR_HMAC_SHA256":
		dekT = aead.AES128CTRHMACSHA256KeyTemplate()
	default:
		log.Fatalf("DEK template %s, is not supported. Expecting AES128_GCM or AES128_CTR_HMAC_SHA256", dek)
	}
	switch strings.ToUpper(kms) {
	case "GCP":
		gcpclient, err := gcpkms.NewGCPClient(gcpURI)
		if err != nil {
			log.Fatal(err)
		}
		_, err = gcpclient.LoadCredentials(gcpCredFile)
		if err != nil {
			log.Fatal(err)
		}
		registry.RegisterKMSClient(gcpclient)
		kh, err = keyset.NewHandle(aead.KMSEnvelopeAEADKeyTemplate(gcpURI, dekT))
		if err != nil {
			log.Fatal(err)
		}
	case "AWS":
		awsclient, err := awskms.NewAWSClient(awsURI)
		if err != nil {
			log.Fatal(err)
		}
		_, err = awsclient.LoadCredentials(awsCredFile)
		if err != nil {
			log.Fatal(err)
		}
		registry.RegisterKMSClient(awsclient)
		kh, err = keyset.NewHandle(aead.KMSEnvelopeAEADKeyTemplate(awsURI, dekT))
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("KMS %s, is not supported. Expecting AWS or GCP", kms)
	}
	ks := insecurecleartextkeyset.KeysetMaterial(kh)
	h, err := insecurecleartextkeyset.Read(&keyset.MemReaderWriter{Keyset: ks})
	if err != nil {
		log.Fatal(err)
	}
	if err := insecurecleartextkeyset.Write(h, keyset.NewBinaryWriter(&b)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(f, b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
