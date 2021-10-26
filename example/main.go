package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/blins/go-cmd"
)

type printCmd struct {
	fs      *flag.FlagSet
	message string
}

func (cmd *printCmd) GetFlags() *flag.FlagSet {
	return cmd.fs
}

func (cmd *printCmd) ParseArgs(args []string) ([]string, error) {
	// no additional aprameters
	return args, nil
}

func (cmd *printCmd) Run(ctx context.Context) (cmd.Waiter, error) {
	fmt.Println(cmd.message)
	return nil, nil
}

func printFabric() cmd.Command {
	cmd := &printCmd{
		fs: flag.NewFlagSet("print", flag.ContinueOnError),
	}
	cmd.fs.StringVar(&cmd.message, "message", "", "message are printed")
	return cmd
}

func main() {

	cmd.RegisterFabric("print", cmd.CommandFabricFunc(printFabric))
	flag.Parse()
	waiter := cmd.ParseAndRun(nil, context.Background())
	waiter.Wait()
}
