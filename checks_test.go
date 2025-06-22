package genki_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oaiiae/genki"
)

func TestChecks_Run(t *testing.T) {
	anError := errors.New("an error")

	testCases := []struct {
		desc      string
		checks    genki.Checks
		targetErr error
	}{
		{
			desc: "nil",
			checks: genki.Checks{
				"a": func(context.Context) error { return nil },
			},
			targetErr: nil,
		},
		{
			desc: "anError",
			checks: genki.Checks{
				"b": func(context.Context) error { return anError },
			},
			targetErr: anError,
		},
		{
			desc: "nil and anError",
			checks: genki.Checks{
				"a": func(context.Context) error { return nil },
				"b": func(context.Context) error { return anError },
			},
			targetErr: anError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.checks.Run(context.Background())
			assert.ErrorIs(t, err, tC.targetErr)
		})
	}
}

func ExampleChecks_Handler_ok() {
	w := httptest.NewRecorder()
	genki.Checks{"cond": func(context.Context) error { return nil }}.
		Handler(func(err error) { fmt.Println("error handler:", err) }).
		ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://example.com", nil))

	fmt.Fprintln(os.Stdout, w.Result().Status)
	fmt.Fprintln(os.Stdout, w.Result().Header)
	io.Copy(os.Stdout, w.Result().Body)

	//Output:
	// 200 OK
	// map[Content-Type:[text/plain; charset=utf-8]]
	// OK
}

func ExampleChecks_Handler_ko() {
	w := httptest.NewRecorder()
	genki.Checks{"cond": func(context.Context) error { return errors.New("unmet") }}.
		Handler(func(err error) { fmt.Println("error handler:", err) }).
		ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://example.com", nil))

	fmt.Fprintln(os.Stdout, w.Result().Status)
	fmt.Fprintln(os.Stdout, w.Result().Header)
	io.Copy(os.Stdout, w.Result().Body)

	//Output:
	// error handler: [cond] unmet
	// 500 Internal Server Error
	// map[Content-Type:[text/plain; charset=utf-8]]
	// [cond] unmet
}

func TestAlways(t *testing.T) {
	checks := genki.Always()
	assert.NoError(t, checks.Run(context.Background()))
}

func TestAfter(t *testing.T) {
	checks := genki.After(time.Millisecond)
	require.Error(t, checks.Run(context.Background()))
	time.Sleep(time.Millisecond)
	assert.NoError(t, checks.Run(context.Background()))
}
