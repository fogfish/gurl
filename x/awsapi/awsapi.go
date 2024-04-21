//
// Copyright (C) 2019 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package awsapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	net "net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/fogfish/gurl/v2/http"
)

// Configure HTTP Stack to use AWS Sign V4
func WithSignatureV4(config aws.Config) http.Config {
	return func(p *http.Protocol) {
		p.Socket = &signer{
			config: config,
			signer: v4.NewSigner(),
			socket: p.Socket,
		}
	}
}

type signer struct {
	config aws.Config
	signer *v4.Signer
	socket http.Socket
}

func (s *signer) Do(req *net.Request) (*net.Response, error) {
	credential, err := s.config.Credentials.Retrieve(req.Context())
	if err != nil {
		return nil, err
	}

	hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if req.Body != nil {
		buf := &bytes.Buffer{}
		hasher := sha256.New()
		stream := io.TeeReader(req.Body, hasher)
		if _, err := io.Copy(buf, stream); err != nil {
			return nil, err
		}
		hash = hex.EncodeToString(hasher.Sum(nil))

		req.Body.Close()
		req.Body = io.NopCloser(buf)
	}

	err = s.signer.SignHTTP(
		req.Context(),
		credential,
		req,
		hash,
		"execute-api",
		s.config.Region,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	return s.socket.Do(req)
}
