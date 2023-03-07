//
// Copyright (C) 2019 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

package http

import (
	"context"
	"encoding/json"
	"io"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/fogfish/gurl/v2"
)

type Status struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Duration string `json:"duration"`
	Reason   string `json:"reason,omitempty"`
	Payload  string `json:"payload"`
}

// Evaluates sequence of tests, returns status object for each
func Once(stack Stack, tests ...func() Arrow) []Status {
	status := make([]Status, len(tests))

	for i, test := range tests {
		arr := test()
		ctx := stack.WithContext(context.Background())

		t := time.Now()
		err := ctx.IO(arr)
		status[i] = newStatus(ctx, arrowName(test), time.Since(t), err)
	}

	return status
}

func WriteOnce(w io.Writer, stack Stack, tests ...func() Arrow) error {
	seq := Once(stack, tests...)

	if bytes, err := json.MarshalIndent(seq, "", "  "); err == nil {
		if _, err := w.Write(bytes); err != nil {
			return err
		}
		return nil
	}

	if _, err := w.Write([]byte{'[', ']'}); err != nil {
		return err
	}
	return nil
}

func newStatus(ctx *Context, id string, dur time.Duration, err error) Status {
	switch v := (err).(type) {
	case nil:
		return Status{
			ID:       id,
			Status:   "success",
			Duration: dur.String(),
			Payload:  string(ctx.Payload),
		}
	case *gurl.NoMatch:
		return Status{
			ID:       id,
			Status:   "failure",
			Duration: dur.String(),
			Reason:   v.Diff,
			Payload:  string(ctx.Payload),
		}
	default:
		return Status{
			ID:       id,
			Status:   "failure",
			Duration: dur.String(),
			Reason:   err.Error(),
			Payload:  string(ctx.Payload),
		}
	}
}

func arrowName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if strings.HasPrefix(name, "main.") {
		name = name[5:]
	}

	return name
}
