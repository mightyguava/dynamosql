package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/xo/usql/drivers"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
	"github.com/xo/usql/text"

	"github.com/mightyguava/dynamosql"
)

func main() {
	var err error

	// load current user
	cur, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	args := NewArgs()

	// run
	err = run(args, cur)
	if err != nil && err != io.EOF && err != rline.ErrInterrupt {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		os.Exit(1)
	}
}

// run processes args, processing args.CommandOrFiles if non-empty, if
// specified, otherwise launch an interactive readline from stdin.
func run(args *Args, u *user.User) error {
	var err error

	// get working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// handle variables
	for _, v := range args.Variables {
		if i := strings.Index(v, "="); i != -1 {
			env.Set(v[:i], v[i+1:])
		} else {
			env.Unset(v)
		}
	}
	for _, v := range args.PVariables {
		if i := strings.Index(v, "="); i != -1 {
			if _, err = env.Pset(v[:i], v[i+1:]); err != nil {
				return err
			}
		} else {
			if _, err = env.Ptoggle(v, ""); err != nil {
				return err
			}
		}
	}

	// create input/output
	l, err := rline.New(len(args.CommandOrFiles) != 0, args.Out, env.HistoryFile(u))
	if err != nil {
		return err
	}
	defer l.Close()

	// create handler
	h := handler.New(l, u, wd, args.NoPassword)

	// open dsn
	config := aws.NewConfig()
	if args.AWSRegion != "" {
		config.Region = &args.AWSRegion
	}
	if args.AWSEndpointURL != "" {
		config.Endpoint = &args.AWSEndpointURL
	}
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	if args.AWSProfile != "" {
		os.Setenv("AWS_PROFILE", args.AWSProfile)
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return err
	}
	driver := dynamosql.New(dynamosql.Config{Session: sess, AlwaysConvertCollectionsToGoType: true})
	sql.Register("dynamosql", driver)
	drivers.Register("dynamosql", drivers.Driver{})
	if err = h.Open("dynamosql", ""); err != nil {
		return err
	}

	// start transaction
	if args.SingleTransaction {
		if h.IO().Interactive() {
			return text.ErrSingleTransactionCannotBeUsedWithInteractiveMode
		}
		if err = h.Begin(); err != nil {
			return err
		}
	}

	// rc file
	if rc := env.RCFile(u); !args.NoRC && rc != "" {
		if err = h.Include(rc, false); err != nil && err != text.ErrNoSuchFileOrDirectory {
			return err
		}
	}

	// setup runner
	f := h.Run
	if len(args.CommandOrFiles) != 0 {
		f = runCommandOrFiles(h, args.CommandOrFiles)
	}

	// run
	if err = f(); err != nil {
		return err
	}

	// commit
	if args.SingleTransaction {
		return h.Commit()
	}

	return nil
}

// runCommandOrFiles proccesses all the supplied commands or files.
func runCommandOrFiles(h *handler.Handler, commandsOrFiles []CommandOrFile) func() error {
	return func() error {
		for _, x := range commandsOrFiles {
			h.SetSingleLineMode(x.Command)
			if x.Command {
				h.Reset([]rune(x.Value))
				if err := h.Run(); err != nil && err != io.EOF {
					return err
				}
			} else {
				if err := h.Include(x.Value, false); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
