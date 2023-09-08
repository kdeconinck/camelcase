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
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kdeconinck/slices"
)

// Holds information about a single rune.
type runeInfo struct {
	r rune
}

// Checks whether or not the rune represented by rInfo is a digit.
func (rInfo *runeInfo) isDigit() bool {
	return unicode.IsDigit(rInfo.r)
}

// Checks whether or not the rune represented by rInfo is an uppercase rune.
func (rInfo *runeInfo) isUppercase() bool {
	return unicode.IsUpper(rInfo.r)
}

// A reader designed for reading "CamelCase" strings.
type rdr struct {
	input       string   // The data this reader operates on.
	pos         int      // The position of this reader.
	hasNextRune bool     // A flag indicating if there's a next rune.
	rdRune      runeInfo // Information about the last rune that was read.
	nxtRune     runeInfo // Information about the next rune that's about to be read.
}

// Read the next rune from r.
func (r *rdr) readRune() {
	r.rdRune = runeInfo{rune(r.input[r.pos])}
	r.pos = r.pos + 1
	r.hasNextRune = r.pos < len(r.input)

	if r.hasNextRune {
		r.nxtRune = runeInfo{rune(r.input[r.pos])}
	}
}

// Undo the last rune from r.
func (r *rdr) unreadRune() {
	r.pos = r.pos - 1
	r.nxtRune = r.rdRune
	r.rdRune = runeInfo{rune(r.input[r.pos])}
	r.hasNextRune = true // NOTE: An undo operation means that there will be always a next rune.
}

// Verify if the word that's currently read by r is a word that should NOT be split.
// If noSplit contains a word that starts with the word that's currently read by r, this function returns true, false
// otherwise.
func (r *rdr) isNoSplitWord(sIdx int, noSplit []string) bool {
	return slices.ContainsFn(noSplit, r.input[sIdx:r.pos+1], func(got, want string) bool {
		return strings.HasPrefix(got, want)
	})
}

// Read the next part from r.
// Each word in noSplit (if provided) is treated as a word that shouldn't be split.
func (r *rdr) readNextPart(noSplit []string) string {
	sIdx := r.pos

	r.readRune()

	if r.rdRune.isDigit() {
		return r.readNumber(sIdx, noSplit)
	}

	return r.readWord(sIdx, noSplit)
}

// Read and return a number from r.
func (r *rdr) readNumber(sIdx int, noSplit []string) string {
	if r.hasNextRune && r.nxtRune.isDigit() {
		for r.hasNextRune && (r.nxtRune.isDigit() || r.isNoSplitWord(sIdx, noSplit)) {
			r.readRune()
		}

		return r.input[sIdx:r.pos]
	}

	return r.input[sIdx:r.pos]
}

// Read and return a word from r.
func (r *rdr) readWord(sIdx int, noSplit []string) string {
	if r.hasNextRune && r.nxtRune.isUppercase() {
		for r.hasNextRune && (r.nxtRune.isUppercase() || r.isNoSplitWord(sIdx, noSplit)) {
			r.readRune()
		}

		if r.hasNextRune && (!r.nxtRune.isUppercase() && !r.nxtRune.isDigit()) {
			r.unreadRune()
		}

		return r.input[sIdx:r.pos]
	}

	for r.hasNextRune && (r.isNoSplitWord(sIdx, noSplit) || (!r.nxtRune.isUppercase() && !r.nxtRune.isDigit())) {
		r.readRune()
	}

	return r.input[sIdx:r.pos]
}

// Split reads v treating it as a "CamelCase" and returns the different words.
// If v isn't a valid UTF-8 string, or when v is an empty string, a slice with one element (v) is returned.
// Each word in noSplit (if provided) is treated as a word that shouldn't be split.
func Split(v string, noSplit ...string) []string {
	if !utf8.ValidString(v) || len(v) == 0 {
		return []string{v}
	}

	vRdr := &rdr{input: v}
	retVal := make([]string, 0)

	for vRdr.pos < len(v) {
		retVal = append(retVal, vRdr.readNextPart(noSplit))
	}

	return retVal
}
