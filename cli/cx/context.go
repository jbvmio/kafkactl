// Copyright © 2018 NAME HERE <jbonds@jbvm.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cx

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

// Context contains parameters for a context.
type Context struct {
	Name          string   `yaml:"name"`
	Brokers       []string `yaml:"brokers"`
	Burrow        []string `yaml:"burrow"`
	Zookeeper     []string `yaml:"zookeeper"`
	ClientVersion string   `yaml:"clientVersion"`
	Ssl           *SSL     `yaml:"ssl"`
	Sasl          *SASL    `yaml:"sasl"`
}

// SASL contains SASL Authentication Parameters.
type SASL struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SSL contains SASL Authentication Parameters.
type SSL struct {
	Insecure bool   `yaml:"insecure"`
	TLSCert  string `yaml:"tlscert"`
	TLSCA    string `yaml:"tlsca"`
	TLSKey   string `yaml:"tlskey"`
}

// SetupCerts takes the paths to a tls certificate, CA, and certificate key in
// a PEM format and returns a constructed tls.Config object.
func SetupCerts(certPath, caPath, keyPath string, insecure bool) (*tls.Config, error) {
	if certPath == "" && caPath == "" && keyPath == "" {
		return nil, nil
	}

	var caPool *x509.CertPool = nil
	if caPath != "" {
		caString, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, err
		}

		caPool = x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM(caString)
		if !ok {
			err := fmt.Errorf("unable to add ca at %s to certificate pool", caPath)
			return nil, err
		}
	}

	clientCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	bundle := &tls.Config{
		RootCAs:            caPool,
		Certificates:       []tls.Certificate{clientCert},
		Renegotiation:      tls.RenegotiateOnceAsClient,
		InsecureSkipVerify: insecure,
	}
	bundle.BuildNameToCertificate()
	return bundle, nil
}
