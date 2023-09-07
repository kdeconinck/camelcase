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

// Quality assurance: Verify (and measure the performance) of the public API of the "camelcase" package.
package camelcase_test

import (
	"strings"
	"testing"

	"github.com/kdeconinck/assert"
	"github.com/kdeconinck/camelcase"
)

// UT: Split a "CamelCase" word into a slice of words.
func TestSplit(t *testing.T) {
	for _, tc := range []struct {
		vInput   string
		vNoSplit []string
		want     []string
	}{
		{
			vInput: "",
			want:   []string{""},
		},
		{
			vInput: "lowercase",
			want:   []string{"lowercase"},
		},
		{
			vInput: "Uppercase",
			want:   []string{"Uppercase"},
		},
		{
			vInput: "MultipleWords",
			want:   []string{"Multiple", "Words"},
		},
		{
			vInput: "HTML",
			want:   []string{"HTML"},
		},
		{
			vInput: "PDFLoader",
			want:   []string{"PDF", "Loader"},
		},
		{
			vInput: "11",
			want:   []string{"11"},
		},
		{
			vInput: "10Validators",
			want:   []string{"10", "Validators"},
		},
		{
			vInput: "GL11Version",
			want:   []string{"GL", "11", "Version"},
		},
		{
			vInput: "5May2000",
			want:   []string{"5", "May", "2000"},
		},
		{
			vInput:   "1Tls2IsUsedInHttpCommunicationAndIsSecure",
			vNoSplit: []string{"Tls2", "HttpCommunication"},
			want:     []string{"1", "Tls2", "Is", "Used", "In", "HttpCommunication", "And", "Is", "Secure"},
		},
		{
			vInput: "BadUTF8\xe2\xe2\xa1",
			want:   []string{"BadUTF8\xe2\xe2\xa1"},
		},
	} {
		// ACT.
		got := camelcase.Split(tc.vInput, tc.vNoSplit...)

		// ASSERT.
		assert.EqualS(t, got, tc.want, "", "\n\n"+
			"UT Name:  Compare 2 slices for equality.\n"+
			"Input:    %v\n"+
			"\033[32mExpected: %v\033[0m\n"+
			"\033[31mActual:   %v\033[0m\n\n", tc.vInput, tc.want, got)
	}
}

// Benchmark: Split a "CamelCase" string.
func BenchmarkSplit(b *testing.B) {
	// ARRANGE.
	var s strings.Builder

	for i := 0; i < 1_000_000; i++ {
		s.WriteString("HelloWorld99HTML")
	}

	input := s.String()

	// RESET.
	b.ResetTimer()

	// EXECUTION.
	for i := 0; i < b.N; i++ {
		_ = camelcase.Split(input)
	}
}
