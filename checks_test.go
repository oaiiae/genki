package genki_test

import (
	"context"
	"fmt"
	"io"
	"time"

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

func ExampleAlways() {
	fmt.Println(genki.Always().Run(context.Background()))

	//Output:
	// <nil>
}

func ExampleAfter() {
	fmt.Println(genki.After(time.Hour).Run(context.Background()))

	//Output:
	// [after] xx left
}
