package gen

import (
	"flag"
	"fmt"
	"os"
)

type cmdgen struct {
	outfile string
}

func (c *cmdgen) Name() string {
	return "gen"
}

func (c *cmdgen) Desc() string {
	return "generate rc skeleton"
}

func (c *cmdgen) Run(args []string) error {
	fset := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fset.StringVar(&c.outfile, "o", "",
		"output file name(default to stdout)")
	err := fset.Parse(args[1:])
	if err != nil {
		return err
	}

	if c.outfile == "" {
		fmt.Println(
			string(MustAsset("skeleton.ank")))
		return nil
	}

	out, err := os.Create(c.outfile)
	if err != nil {
		return err
	}
	defer out.Close()

	fmt.Fprintln(out, string(MustAsset("skeleton.ank")))

	return nil
}

func New() *cmdgen {
	return &cmdgen{}
}
