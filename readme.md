# MuxChain

MuxChain is a small package designed to complement [net/http](http://golang.org/pkg/net/http) for specifying chains of handlers. With it, you can succinctly compose layers of middleware without introducing large dependencies or effectively defeating the type system.

#### Example

```
muxchain.Chain("/", logger, gzipHandler, echoHandler)
http.ListenAndServe(":36363", muxchain.Default)
```

This specifies that all patterns matched should be handled by the logger, then gzip, then echo. Since we're chaining to the default MuxChain, we can just pass that to `http.ListenAndServe`. You can see a more complete example in the "sample" directory.

### Status and Direction

The naming or APIs may change a little bit to make code using this package a bit more readable. The package is also quite new, so ideas may emerge for improvement after some usage. Please feel free to make these types of suggestions, particularly in this early stage.

#### License

BSD 3-clause (see LICENSE file)