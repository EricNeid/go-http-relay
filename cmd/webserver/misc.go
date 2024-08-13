// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT
package main

import "io"

type lazyMultiWriter struct {
	writers []io.Writer
}

func (t *lazyMultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		w.Write(p)
	}
	return len(p), nil
}

var _ io.StringWriter = (*lazyMultiWriter)(nil)

func (t *lazyMultiWriter) WriteString(s string) (n int, err error) {
	var p []byte // lazily initialized if/when needed
	for _, w := range t.writers {
		if sw, ok := w.(io.StringWriter); ok {
			sw.WriteString(s)
		} else {
			if p == nil {
				p = []byte(s)
			}
			w.Write(p)
		}
	}
	return len(s), nil
}

// LazyMultiWriter creates a writer that duplicates its writes to all the
// provided writers. It is very similar to io.MultiWriter, but it does not
// stop writing if any write operation fails. It simply continues down the list.
func LazyMultiWriter(writers ...io.Writer) io.Writer {
	io.MultiWriter()
	allWriters := make([]io.Writer, 0, len(writers))
	for _, w := range writers {
		if mw, ok := w.(*lazyMultiWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &lazyMultiWriter{allWriters}
}
