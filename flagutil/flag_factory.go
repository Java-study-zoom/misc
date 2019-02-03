package flagutil

import (
	"flag"
)

// FlagFactory is a factory that can be used to create flag sets.
type FlagFactory struct {
	name string
}

// ParseArgs parses the args using the given FlagSet, and returns
// the parsed args that does not include flags.
func ParseArgs(s *flag.FlagSet, args []string) []string {
	s.Parse(args)
	return s.Args()
}

// Make creates a new flag set.
func (f *FlagFactory) Make() *flag.FlagSet {
	return flag.NewFlagSet(f.name, flag.ExitOnError)
}

// PlainArgs parse the args with no flags.
func (f *FlagFactory) PlainArgs(args []string) []string {
	flags := f.Make()
	return ParseArgs(flags, args)
}
