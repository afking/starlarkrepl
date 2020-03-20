package starlarkrepl

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/peterh/liner"
	"go.starlark.net/repl"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func Run(thread *starlark.Thread, globals starlark.StringDict) (err error) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	resolve.LoadBindsGlobally = true // TODO

	for err == nil {
		err = rep(line, thread, globals)
	}
	if err == io.EOF {
		fmt.Println()
		return nil
	}
	return err
}

// finxPrefix assumes sorted arrays of keys
func findPrefix(pfx string, keyss ...[]string) (c []string) {
	for _, keys := range keyss {
		i := sort.SearchStrings(keys, pfx)
		j := i
		for ; j < len(keys); j++ {
			if !strings.HasPrefix(keys[j], pfx) {
				break
			}
		}
		c = append(c, keys[i:j]...)
	}
	sort.Strings(c)
	return c
}

func newCompleter(globals starlark.StringDict) liner.Completer {
	return func(line string) (c []string) {
		if strings.Count(line, " ") == len(line) {
			return []string{strings.Repeat(" ", (len(line)/4)*4+4)}
		}
		var wrote bool
		f, err := syntax.ParseCompoundStmt("<stdin>", func() ([]byte, error) {
			if wrote {
				return nil, io.EOF
			}
			wrote = true
			return []byte(line + "\n"), nil
		})
		if err != nil {
			return
		}
		fmt.Println(f, err)

		if len(f.Stmts) == 0 {
			return nil
		}
		var predeclared starlark.StringDict
		if err := resolve.REPLChunk(f, globals.Has, predeclared.Has, starlark.Universe.Has); err != nil {
			//fmt.Println("err", err)
			//return []string{err.Error()}
			// igore name pass failures
		}
		syntax.Walk(f.Stmts[0], func(n syntax.Node) bool {
			if n == nil {
				return false
			}
			start, end := n.Span()
			fmt.Println(start.Col, end.Col)

			switch n := n.(type) {
			case *syntax.Ident:
				fmt.Printf("%T\t%+v\n", n, n)
				fmt.Println("\t---", n.Name, n.Binding)
				keys := globals.Keys()
				c = findPrefix(n.Name, keys, starlark.Universe.Keys())

			case *syntax.DotExpr:
				fmt.Printf("%T\t%+v\n", n, n)
				fmt.Println("\t---", n.X)
				fmt.Println("\t---", n.Name)

			default:
				fmt.Printf("%T\t%+v\n", n, n)
			}
			return true
		})

		//locals, err := resolve.Expr(f, globals.Has, starlark.Universe.Has)
		//if err != nil {
		//	fmt.Println("err", err)
		//}
		//fmt.Prinln("local", locals)

		return
	}
}

func suggest(line string) string {
	var noSpaces int
	for _, c := range line {
		if c == ' ' {
			noSpaces += 1
		} else {
			break
		}
	}
	if strings.HasSuffix(line, ":") {
		noSpaces += 4
	}
	return strings.Repeat(" ", noSpaces)
}

func rep(line *liner.State, thread *starlark.Thread, globals starlark.StringDict) error {
	ctx := context.Background()
	thread.SetLocal("context", ctx)

	var eof bool
	var previous string
	prompt := ">>> "
	readline := func() ([]byte, error) {
		text := suggest(previous)
		s, err := line.PromptWithSuggestion(prompt, text, -1)
		if err != nil {
			switch err {
			case io.EOF:
				eof = true
			case liner.ErrPromptAborted:
				return []byte("\n"), nil
			}
			return nil, err
		}
		prompt = "... "
		previous = s
		return []byte(s + "\n"), nil
	}

	//line.SetCompleter(newCompleter(globals))

	f, err := syntax.ParseCompoundStmt("<stdin>", readline)
	if err != nil {
		if eof {
			return io.EOF
		}
		repl.PrintError(err)
		return nil
	}

	if expr := soleExpr(f); expr != nil {
		// eval
		v, err := starlark.EvalExpr(thread, expr, globals)
		if err != nil {
			repl.PrintError(err)
			return nil
		}

		// print
		if v != starlark.None {
			fmt.Println(v)
		}
	} else if err := starlark.ExecREPLChunk(f, thread, globals); err != nil {
		repl.PrintError(err)
		return nil
	}
	return nil
}

func soleExpr(f *syntax.File) syntax.Expr {
	if len(f.Stmts) == 1 {
		if stmt, ok := f.Stmts[0].(*syntax.ExprStmt); ok {
			return stmt.X
		}
	}
	return nil
}
