# memless — Go bridge

A Go bridge to the memless C ABI through [purego](https://github.com/ebitengine/purego)
— no cgo, no build step besides the pure-Go module. It loads the same
`libmemless_capi` cdylib as the PHP and Node bridges and speaks the same
contract (ABI version 5), so the three stay in parity from a single shared
surface. See the
[project README](https://github.com/nascent-tech/memless#readme) for what
memless is, the guessing rules and the supported SQL subset.

## Install

```sh
go get github.com/nascent-tech/memless-go
```

That is all: the module carries the native library, so there is nothing to
download or configure. Go 1.21 or later, on macOS (`darwin-arm64`,
`darwin-x64`) or Linux with glibc (`linux-x64-gnu`, `linux-arm64-gnu`);
elsewhere, see [The cdylib](#the-cdylib).

The module is published from `bindings/go` of
[nascent-tech/memless](https://github.com/nascent-tech/memless) into the
mirror repository `nascent-tech/memless-go`, which the release workflow fills
at each version with these sources and the four libraries under
`lib/<platform>/` and their `lib/SHA256SUMS`, as a single commit tagged
`vX.Y.Z`. The main repository carries no binary. Your program embeds the
library of the platform it is built for, and only that one (about 5 MB).

Versions up to 0.2.1 were published at
`github.com/nascent-tech/memless/bindings/go` (tags `bindings/go/vX.Y.Z` of
the main repository) and stay there; from 0.3.0 on, the module path is
`github.com/nascent-tech/memless-go`. To move, replace the import path and
run `go get github.com/nascent-tech/memless-go`.

## Surface

```go
package main

import (
	"fmt"
	"log"

	memless "github.com/nascent-tech/memless-go"
)

func main() {
	db, err := memless.Load("data.yaml")
	if err != nil {
		log.Fatal(err) // *memless.RefusalError or *memless.FaultError
	}
	defer db.Release()

	rows, err := db.Query("SELECT name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows.Columns, rows.Rows) // [name] [[Ada] [Grace]]

	affected, err := db.Execute("UPDATE users SET name = 'Zoe' WHERE id = '01H7B2'")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(affected)

	if err := db.Begin(); err != nil { // BEGIN
		log.Fatal(err)
	}
	if _, err := db.Execute("DELETE FROM wallets WHERE id = 'w_123'"); err != nil {
		log.Fatal(err)
	}
	if err := db.Commit(); err != nil { // COMMIT
		log.Fatal(err)
	}

	if err := db.Reload(); err != nil { // re-reads data.yaml
		log.Fatal(err)
	}
}
```

- `Instance.Query` returns `Rows{Columns []string, Rows [][]any}`; a cell is
  `int64`, `float64`, `string`, `bool`, or `nil` for an absent value.
- `Instance.Execute` returns the affected row count (`0` for `BEGIN` /
  `COMMIT` / `ROLLBACK`, which are also run through it — `Begin`, `Commit` and
  `Rollback` are thin wrappers over it).
- Every method returns a `*RefusalError` (the domain's message, verbatim —
  D13) or a `*FaultError` (a boundary or internal fault); tell them apart with
  `errors.As`. A fault reads `memless fault (<status>): <message>`, the same
  text as in the PHP and Node bridges, and carries the ABI status in
  `FaultError.Status` (`2` for an invalid argument, such as a released
  instance or a NUL byte in the SQL, `3` for an internal fault).
- A library that cannot be found or speaks another ABI version is a plain
  error, returned by the first call that needs the library.
- `Instance.Release` releases the handle once; a second call, a zero-value
  instance, or an unknown handle is ignored.
- `Instance.Reload` re-reads the file from disk into a fresh in-memory state,
  exactly as `Load` would build it. It returns a `*RefusalError` while a
  transaction is open (`cannot reload while a transaction is open`) or when
  the file would be refused by `Load`; the old state stays usable either way.

## The cdylib

The bridge loads the native library on the first call that needs it, never
at import time, and looks for it in the same order as the PHP and Node
bridges:

1. `MEMLESS_LIB`, a trusted (ideally absolute) path that must name an
   existing file;
2. the library embedded for this platform, extracted once into
   `<user cache dir>/memless/<version>-<short sha256>/` (`~/Library/Caches` on
   macOS, `$XDG_CACHE_HOME` or `~/.cache` on Linux). The directory is created
   `0700`, the file is written to a temporary name then renamed, and before
   every load its full SHA-256 is compared with the embedded bytes and the file
   rewritten when they differ. On Linux the bridge checks the libc
   (`/usr/bin/ldd`, then the dynamic loader under `/lib`), because the
   bundled libraries need glibc: musl, or a libc it cannot tell, has none;
3. inside a checked-out workspace, `target/release/`, then `target/debug/`
   (`.dylib` before `.so`).

When the cache directory cannot be used, the bridge moves on to step 3. If the
user cache directory is mounted `noexec`, the extracted library cannot be
loaded: set `MEMLESS_LIB`. The bundled Linux libraries need glibc 2.39 or
later (Ubuntu 24.04 or later); on an older glibc, set `MEMLESS_LIB` to a
library built locally. When nothing is found, the error says to set
`MEMLESS_LIB`. The library must speak ABI version 5. On another platform, or
with your own build, set `MEMLESS_LIB`:

```sh
cargo build --release -p memless-capi
export MEMLESS_LIB="$PWD/target/release/libmemless_capi.so"   # .dylib on macOS
```

`MEMLESS_LIB` loads arbitrary native code, like any FFI library path — only
point it at a library you trust.

## Running the tests

The tests live in `bindings/go` of the main repository,
[nascent-tech/memless](https://github.com/nascent-tech/memless); the mirror
does not carry them:

```sh
cargo build -p memless-capi
cd bindings/go && go test ./...
```

## Contributing

This bridge is developed in [`bindings/go`](https://github.com/nascent-tech/memless/tree/main/bindings/go)
of [nascent-tech/memless](https://github.com/nascent-tech/memless), next to
the Rust core, the two other bridges and the parity harness that keeps the
three in agreement. The repository `nascent-tech/memless-go` is a mirror
that the release workflow rewrites at each version: open issues and pull
requests on [nascent-tech/memless](https://github.com/nascent-tech/memless/issues/new/choose),
and read its [contributing guide](https://github.com/nascent-tech/memless/blob/main/CONTRIBUTING.md)
first. Security issues go to its
[security policy](https://github.com/nascent-tech/memless/security/policy).
