//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//
func encode(content string, data interface{}) (buf *bytes.Buffer, err error) {
	switch content {
	case "application/json":
		buf, err = encodeJSON(data)
	// case "application/x-www-form-urlencoded":
	// 	req.payload, req.fail = encodeForm(data)
	default:
		err = fmt.Errorf("unsupported Content-Type %v", content)
	}

	return
}

func encodeJSON(data interface{}) (*bytes.Buffer, error) {
	json, err := json.Marshal(data)
	return bytes.NewBuffer(json), err
}
