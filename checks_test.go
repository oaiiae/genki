package genki_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

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
