package genki_test

import (
	"context"
	"fmt"
	"io"

	"github.com/oaiiae/genki"
)

func ExampleChecks_Run() {
	fmt.Println(genki.Checks{
		"a": func(context.Context) error { return nil },
		"b": func(context.Context) error { return io.EOF },
	}.Run(context.Background()))

	//Output:
	// [b] EOF
}
