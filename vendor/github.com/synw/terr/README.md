# Traced errors

Error tracing library

## Usage

   ```go
package main

import (
	"errors"
	"github.com/synw/terr"
)

func f2() *terr.Trace {
	err := errors.New("error from f2")
	newerr := terr.Debug("f2", err)
	return newerr
}

func f1() *terr.Trace {
	err := errors.New("error from f1")
	perr := f2()
	newerr := terr.New("f1", err, perr)
	return newerr
}

func main() {
	err := f1()
	if err != nil {
		err.Print()
	}
}
```

Start tracing errors:

   ```go
// terr.New(from string, err error)
trace := terr.New("function_path", err)
return trace
//return a *terr.Trace instead of an error
   ```

Continue tracing:

   ```go
// trace is the previous returned *terr.Trace
terr.Add("function_path", err, trace)
// pass the trace without adding a new error
terr.Pass("function_path", trace)
   ```

## Options

   ```go
terr.Critical("function_path", err)
terr.Minor("function_path", err)
terr.Debug("function_path", err)
   ```
   
Print the errors as they come:

   ```go
terr.Push("function_path", err, previous_trace)
// print the go stack trace as well
terr.Stack("function_path", err, previous_trace)
   ```

## Formating

Custom formating:
   ```go
// trace is a *terr.Trace
trace.Print()
// with colors
trace.Printc()
// with prefix and suffix trace.Printps(prefix, suffix)
trace.Printps("->", "\n")
// get the trace output without printing
formated_trace := trace.Format()
// with error class labels
formated_trace := trace.Formatl()
   ```
   
Check the [examples](https://github.com/synw/terr/tree/master/example)
