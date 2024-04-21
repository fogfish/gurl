//
// Copyright (C) 2019 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/gurl
//

// Package xhtml is an extension to gurl library for fetching and parsing
// html using query selectors.
package xhtml

import (
	"bytes"
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/fogfish/gurl/v2/http"
	ø "github.com/fogfish/gurl/v2/http/recv"
	ƒ "github.com/fogfish/gurl/v2/http/send"
)

// HTTP Arrow receive and parses HTML content
func Content(content **goquery.Selection) http.Arrow {
	b := &bytes.Buffer{}

	return http.Join(
		ø.Bytes(b),
		func(ctx *http.Context) error {
			doc, err := goquery.NewDocumentFromReader(b)
			if err != nil {
				return err
			}

			c := doc.Selection
			*content = c

			return nil
		},
	)
}

// HTTP Client tailored for fetching HTML content from public resources
type Client struct {
	http.Stack
}

// Create new instance of HTTP Client tailored for fetching HTML content
// from public resources.
func New(opts ...http.Config) Client {
	return Client{Stack: http.New(opts...)}
}

// Fetches HTML content behind url, parser HTML.
func (site Client) Fetch(ctx context.Context, url string) (*goquery.Selection, error) {
	var content *goquery.Selection

	err := site.Stack.IO(context.Background(),
		http.GET(
			ƒ.URI(url),
			ƒ.Accept.HTML,

			ø.Status.OK,
			Content(&content),
		),
	)
	if err != nil {
		return nil, err
	}

	return content, nil
}
