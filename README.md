# Memless for Go

Query and change a YAML file with SQL, from Go. Memless loads the file into
memory, runs your SQL against it and writes every accepted change back to the
file — no schema, no server. It is made for test fixtures and demos.

Everything you need to use it from Go is on this page. The
[project README](https://github.com/nascent-tech/memless#readme) has more
detail on the file format, transactions and errors.

## Install

```sh
go get github.com/nascent-tech/memless-go
```

Requires Go 1.21 or later, and **no cgo**: the module calls the native engine
through [purego](https://github.com/ebitengine/purego). It embeds the engine
for macOS (Apple silicon, Intel) and Linux with glibc 2.39 or later (x86_64,
aarch64); your program carries only the one for the platform you build for
(about 5 MB).

## Quick start

```yaml
# data.yaml
users:
  - id: 1
    name: Ada
  - id: 2
    name: Grace
```

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
		log.Fatal(err)
	}
	defer db.Release()

	// data.yaml is rewritten before Execute returns.
	if _, err := db.Execute("INSERT INTO users (id, name) VALUES (3, 'Linus')"); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT name FROM users ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows.Columns, rows.Rows) // [name] [[Ada] [Grace] [Linus]]
}
```

## Writing the YAML file

Each top-level key is a table, each table is a list of rows, and each row is a
map of column names to values. Memless reads everything else from the data:

| Rule | What Memless expects |
| --- | --- |
| Row identity | Every row has an `id`, text or integer, unique within its table. `5` and `"5"` are different ids. |
| Relations | A column named `<name>_id` points at the `id` of the table `<name>s`, when it exists: `user_id` → `users`. The plural is always `<name>` + `s` (`category_id` → `categorys`). Every value present must name an existing row; a missing value is allowed. |
| Values | Text, integer, decimal or boolean. Lists and maps inside a row are refused. |
| Missing values | Any column but `id` may be left out; `null` (or `~`) counts as missing and reads back as `nil`. |
| Types never convert | An integer (`2`) and a decimal (`2.0`) are different types. A comparison only matches the same type: `WHERE amount > 1` skips `1.5`, and `WHERE amount = 2.0` does not match `2`. Keep one type per column (`2.0`, not `2`, in a decimal column). |
| Mixed types | Allowed in a column, but sorting or summing rows of different types is refused. |

## Supported SQL

| Statement | Supported form |
| --- | --- |
| `SELECT` | `SELECT <columns> \| * FROM <table>`, with optional `WHERE`, `ORDER BY` and one `JOIN` |
| `WHERE` | `=`, `<>`, `<`, `<=`, `>`, `>=`, `IS NULL`, `IS NOT NULL`, with `AND`, `OR` and parentheses |
| `JOIN` | One `INNER JOIN <table> ON <a>.<name>_id = <b>.id`; every column is then written `table.column` |
| `ORDER BY` | One or more columns, `ASC` (default) or `DESC`; ties keep the file order, missing values come last |
| Aggregates | `SELECT COUNT(*)` or `SELECT SUM(<column>)`, alone in the select list |
| `INSERT` | `INSERT INTO <table> (<columns>) VALUES (<values>)`; the column list is required |
| `UPDATE` / `DELETE` | `UPDATE <table> SET <column> = <value>, … [WHERE …]`, `DELETE FROM <table> [WHERE …]` |

Anything else — `GROUP BY`, `LIMIT`, several joins, `ORDER BY 1`, `LIKE`, `IN`, subqueries,
`CREATE`… — is refused with `<construct> is outside the supported SQL subset`.

## Good to know

- **Writes are saved at once.** Outside a transaction, each accepted write
  rewrites the file before the call returns. The file is replaced atomically
  through a temporary `.<name>.memless-tmp` next to it: add `.*.memless-tmp`
  to your `.gitignore`.
- **In tests, load a copy of your fixture**, or wrap the test in a
  transaction that you roll back, so that every test starts from the same
  data.
- **One instance per file.** Memless does not coordinate writers: do not
  change the same file from two instances or two processes at once.

## API

| Call | Returns | What it does |
| --- | --- | --- |
| `memless.Load(path)` | `*Instance, error` | Reads the YAML file into memory. |
| `db.Query(sql)` | `Rows, error` | Runs a `SELECT`. `Rows.Columns` holds the names, `Rows.Rows` the values (`[][]any`). |
| `db.Execute(sql)` | `uint64, error` | Runs an `INSERT`, `UPDATE` or `DELETE` and returns the number of rows affected. Outside a transaction, the file is rewritten before it returns. |
| `db.Begin()` / `db.Commit()` / `db.Rollback()` | `error` | Groups writes: nothing is written until `Commit`; `Rollback` discards them. |
| `db.Reload()` | `error` | Rereads the file, after something else changed it. Refused while a transaction is open. |
| `db.Release()` | — | Frees the instance. Calling it again does nothing. |

## Values

| In the file | In Go (`any`) |
| --- | --- |
| text | `string` |
| integer | `int64` |
| decimal | `float64` |
| boolean | `bool` |
| missing value | `nil` |

## Errors

Every call returns one of two error types. Their messages are identical in the
Node.js and PHP versions of Memless.

- **`*memless.RefusalError`** — the file or the SQL breaks a rule: unknown
  table or column, broken relation, unsupported SQL, transaction already open.
  It is an expected outcome, which you can check for in a test.
- **`*memless.FaultError`** — the package was misused (for example, a call after
  `Release`) or an internal error happened. Its `Status` field is `2` for an
  invalid argument and `3` for an internal error.

```go
_, err := db.Query("SELECT nope FROM users")

var refusal *memless.RefusalError
if errors.As(err, &refusal) {
	fmt.Println(refusal) // no column "nope" in table "users"
}
```

If the native library cannot be found or has an incompatible version, the
first call returns a plain error that explains what to do.

## Troubleshooting

**Your platform is not covered** (Alpine and other musl-based Linux, glibc
older than 2.39). Build the engine and set `MEMLESS_LIB`:

```sh
git clone https://github.com/nascent-tech/memless.git && cd memless
cargo build --release -p memless-capi
export MEMLESS_LIB="$PWD/target/release/libmemless_capi.so"   # .dylib on macOS
```

**Where the embedded engine goes.** On first use, the engine is written once
to your user cache directory (`~/Library/Caches/memless/` on macOS,
`$XDG_CACHE_HOME/memless/` or `~/.cache/memless/` on Linux) and its SHA-256 is
checked before every load. If that directory cannot be used, the first call
fails with an error naming `MEMLESS_LIB`; if it is mounted `noexec`, loading
fails. In both cases, set `MEMLESS_LIB`.

The package looks for the engine in this order: `MEMLESS_LIB`, then the
embedded engine for your platform, then, inside a clone of the repository,
`target/release/` and `target/debug/`. On Linux it checks the libc first
(`/usr/bin/ldd`, then the dynamic loader under `/lib`): musl, or a libc it
cannot identify, gets no embedded engine. `MEMLESS_LIB` loads native code into
your process; only point it at a library you trust.

**Upgrading from 0.2.x.** The module used to be
`github.com/nascent-tech/memless/bindings/go`. Replace that import path with
`github.com/nascent-tech/memless-go`; the package name, `memless`, is the same.

## Contributing

This module is developed in
[`bindings/go`](https://github.com/nascent-tech/memless/tree/main/bindings/go)
of [nascent-tech/memless](https://github.com/nascent-tech/memless), with the
engine and the Node.js and PHP versions. Open issues and pull requests there —
the repository `nascent-tech/memless-go`, which the Go module proxy reads, is a
read-only copy published with each release. See the
[contributing guide](https://github.com/nascent-tech/memless/blob/main/CONTRIBUTING.md)
and the [security policy](https://github.com/nascent-tech/memless/security/policy).

To work on this package in a clone of the repository:

```sh
cargo build --release -p memless-capi   # in a clone, the package loads target/release/
cd bindings/go && go test ./...
```
