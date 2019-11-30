# maze

Utilities for creating, rendering, and solving mazes.

The core of this library is a Go implementation of the algorithms described in [Mazes for Programmers](https://pragprog.com/book/jbmaze/mazes-for-programmers) by Jamis Buck.

Many of the maze-creation algorithms utilize random selection; for best results, be sure to initialize the default random source before using them:

```go
rand.Seed(time.Now().UnixNano())
```
