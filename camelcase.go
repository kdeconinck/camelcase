// =====================================================================================================================
// = LICENSE:       Copyright (c) 2023 Kevin De Coninck
// =
// =                Permission is hereby granted, free of charge, to any person
// =                obtaining a copy of this software and associated documentation
// =                files (the "Software"), to deal in the Software without
// =                restriction, including without limitation the rights to use,
// =                copy, modify, merge, publish, distribute, sublicense, and/or sell
// =                copies of the Software, and to permit persons to whom the
// =                Software is furnished to do so, subject to the following
// =                conditions:
// =
// =                The above copyright notice and this permission notice shall be
// =                included in all copies or substantial portions of the Software.
// =
// =                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// =                EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// =                OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// =                NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// =                HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// =                WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// =                FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// =                OTHER DEALINGS IN THE SOFTWARE.
// =====================================================================================================================

// Package camelcase defines functions for working with "CamelCase" strings.
package camelcase

import (
	"unicode"
	"unicode/utf8"
)

// A reader designed for reading "CamelCase" strings.
type rdr struct {
	input  string // The data this reader operates on.
	pos    int    // The position of this reader.
	rdRune rune   // The last rune that was read.
}

// Read the next rune from r.
func (r *rdr) readRune() {
	r.rdRune = rune(r.input[r.pos])
	r.pos = r.pos + 1
}

// Undo the last rune from r.
func (r *rdr) unreadRune() {
	r.pos = r.pos - 1
}

// Peek at the next rune from r without advancing the reader.
func (r *rdr) peekRune() rune {
	return rune(r.input[r.pos])
}

// Read the current word from r.
// The word is considered terminated as soon as the reader encounters a new uppercase character.
func (r *rdr) readWord() string {
	sIdx := r.pos

	r.readRune()

	if unicode.IsDigit(r.rdRune) {
		if r.pos < len(r.input) && unicode.IsDigit(r.peekRune()) {
			for r.pos < len(r.input) && unicode.IsDigit(r.peekRune()) {
				r.readRune()
			}

			return r.input[sIdx:r.pos]
		}

		return r.input[sIdx:r.pos]
	} else {
		if r.pos < len(r.input) && (unicode.IsUpper(r.peekRune())) {
			for r.pos < len(r.input) && (unicode.IsUpper(r.peekRune())) {
				r.readRune()
			}

			if r.pos < len(r.input) && (!unicode.IsUpper(r.peekRune()) && !unicode.IsDigit(r.peekRune())) {
				r.unreadRune()
			}

			return r.input[sIdx:r.pos]
		}

		for r.pos < len(r.input) && (!unicode.IsUpper(r.peekRune()) && !unicode.IsDigit(r.peekRune())) {
			r.readRune()
		}

		return r.input[sIdx:r.pos]
	}
}

// Split reads v treating it as a "CamelCase" and returns the different words.
// If v isn't a valid UTF-8 string, or when v is an empty string, a slice with one element (v) is returned.
func Split(v string) []string {
	if !utf8.ValidString(v) || len(v) == 0 {
		return []string{v}
	}

	vRdr := &rdr{input: v}
	retVal := make([]string, 0)

	for vRdr.pos < len(v) {
		retVal = append(retVal, vRdr.readWord())
	}

	return retVal
}
