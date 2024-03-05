// Copyright (c) 2023–present Bartłomiej Krukowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package grouperror_test

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/gontainer/grouperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:goerr113
func TestPrefixedGroup(t *testing.T) {
	t.Parallel()

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, grouperror.Join(nil, nil))
	})

	t.Run("One-error collection", func(t *testing.T) {
		t.Parallel()

		errs := grouperror.Collection(io.EOF)
		assert.Equal(t, []error{io.EOF}, errs)
	})

	t.Run("Errors", func(t *testing.T) {
		t.Parallel()

		fieldsErr := grouperror.Join(
			errors.New("invalid value of `name`"),
			errors.New("invalid value of `age`"),
		)

		personErr := grouperror.Prefix(
			"Person: ",
			fieldsErr,
			errors.New("the given ID does not exist"),
		)
		validationErr := grouperror.Prefix(
			"Validation: ",
			personErr,
			errors.New("unexpected error"),
		)

		errs := grouperror.Collection(validationErr)
		expected := []string{
			"Validation: Person: invalid value of `name`",
			"Validation: Person: invalid value of `age`",
			"Validation: Person: the given ID does not exist",
			"Validation: unexpected error",
		}

		require.Len(t, errs, len(expected), validationErr)
		for i, err := range errs {
			assert.EqualError(t, err, expected[i]) //nolint:testifylint
		}

		assert.EqualError(t, validationErr, strings.Join(expected, "\n")) //nolint:testifylint
	})
}

type wrappedError struct {
	error
}

func (w *wrappedError) Collection() []error {
	return []error{w.error}
}

func TestCollection(t *testing.T) {
	t.Parallel()

	t.Run("Nil", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, grouperror.Collection(nil))
	})

	//nolint:goerr113
	t.Run("Custom error", func(t *testing.T) {
		t.Parallel()

		t.Run("Implements interface{ Collection() []error }", func(t *testing.T) {
			t.Parallel()

			err := &wrappedError{
				error: grouperror.Prefix("my group: ", errors.New("error #1"), errors.New("error #2")),
			}

			expected := []string{
				"my group: error #1",
				"my group: error #2",
			}

			collection := grouperror.Collection(err)
			require.Len(t, collection, 2)

			for i, x := range collection {
				require.EqualError(t, x, expected[i])
			}
		})

		t.Run("Does not implement interface{ Collection() []error }", func(t *testing.T) {
			t.Parallel()

			parent := grouperror.Prefix("my group: ", errors.New("error #1"), errors.New("error #2"))
			err := &wrappedError{
				error: fmt.Errorf("error: %w", parent),
			}

			expected := "error: my group: error #1\n" +
				"my group: error #2"

			collection := grouperror.Collection(err)
			require.Len(t, collection, 1)
			require.EqualError(t, collection[0], expected)
		})
	})
}

func Test_groupError_Unwrap(t *testing.T) {
	t.Parallel()

	const wrongFileName = "file does not exist"

	getPathError := func() error {
		_, err := os.Open(wrongFileName)

		return err //nolint:wrapcheck
	}

	err := grouperror.Prefix(
		"my group: ",
		grouperror.Prefix("some errors: ", io.EOF, io.ErrNoProgress),
		io.ErrUnexpectedEOF,
		getPathError(),
	)

	err = grouperror.Prefix("errors: ", err)

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()
		for _, target := range []error{io.EOF, io.ErrNoProgress, io.ErrUnexpectedEOF} {
			assert.ErrorIs(t, err, target) //nolint:testifylint
		}
		assert.NotErrorIs(t, err, io.ErrClosedPipe) //nolint:testifylint
	})

	t.Run("errors.As", func(t *testing.T) {
		t.Run("*os.PathError", func(t *testing.T) {
			t.Parallel()

			var target *os.PathError
			if assert.ErrorAs(t, err, &target) { //nolint:testifylint
				assert.Equal(t, wrongFileName, target.Path)
			}
		})
		t.Run("*net.AddrError", func(t *testing.T) {
			t.Parallel()

			t.Run("false", func(t *testing.T) {
				t.Parallel()

				var target *net.AddrError
				assert.False(t, errors.As(err, &target))
			})

			t.Run("true", func(t *testing.T) {
				t.Parallel()

				ip := net.IP{1, 2, 3}
				_, addrErr := ip.MarshalText() // address 010203: invalid IP address
				var target *net.AddrError
				assert.Nil(t, target)
				require.ErrorAs(
					t,
					grouperror.Join(err, addrErr),
					&target,
				)
				require.EqualError(t, target, "address 010203: invalid IP address")
			})
		})
	})
}
