//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package gurl

import (
	"encoding/json"
	"fmt"
)

// Status contains summary about the IO
type Status struct {
	ID       string      `json:"id"`
	Status   string      `json:"status"`
	Duration int64       `json:"duration"`
	Expect   interface{} `json:"expect,omitempty"`
	Actual   interface{} `json:"actual,omitempty"`
}

// Status returns the status of IOCat
func (io *IOCat) Status(id string) Status {
	switch v := io.Fail.(type) {
	case nil:
		return Status{
			ID:       id,
			Status:   "success",
			Duration: io.dur.Milliseconds(),
			Actual:   io.Body,
		}
	case *BadMatch:
		return Status{
			ID:       id,
			Status:   "failure",
			Duration: io.dur.Milliseconds(),
			Actual:   v.Actual,
			Expect:   v.Expect,
		}
	default:
		return Status{
			ID:       id,
			Status:   "failure",
			Duration: io.dur.Milliseconds(),
			Actual:   fmt.Sprint(io.Fail),
		}
	}
}

// Tagged is an alias for Arrow type
type Tagged struct {
	Label string
	Arrow func() Arrow
}

// Once evaluates set of tagged arrows
func Once(tagged ...Tagged) []byte {
	status := []Status{}
	for _, f := range tagged {
		arrow := f.Arrow()
		status = append(status, arrow(IO()).Status(f.Label))
	}
	if bytes, err := json.Marshal(status); err == nil {
		return bytes
	}

	return []byte{'{', '}'}
}

// Println evaluates set of tagged arrows and output results
func Println(tagged ...Tagged) {
	fmt.Println(string(Once(tagged...)))
}
