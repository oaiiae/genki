package genki

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Checks map[string]func(context.Context) error

func (m Checks) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for name, check := range m {
		eg.Go(func() error {
			err := check(ctx)
			if err != nil {
				err = fmt.Errorf("[%s] %w", name, err)
			}
			return err
		})
	}
	return eg.Wait()
}

func (m Checks) Handler(errorHandler func(error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		switch err := m.Run(r.Context()); err {
		case nil:
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			errorHandler(err)
		}
	})
}

func Always() Checks {
	return Checks{"always": func(context.Context) error { return nil }}
}

func After(d time.Duration) Checks {
	t := time.Now().Add(d)
	return Checks{"after": func(context.Context) error {
		left := time.Until(t)
		if left > 0 {
			return fmt.Errorf("%v left", left)
		}
		return nil
	}}
}
