# starlarkrepl

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)](https://pkg.go.dev/github.com/afking/starlarkrepl?tab=doc)

Experimental autocompletion REPL for starlark-go based on [liner](https://github.com/peterh/liner).

```go
thread := &starlark.Thread{Load: repl.MakeLoad()}
globals := make(starlark.StringDict)
options := starlarkrepl.Options{
	AutoComplete: true, 
	HistoryFile: "history.txt",
}

starlarkrepl.Run(thread, globals, options)
```
