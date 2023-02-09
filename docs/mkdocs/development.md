### Development

You should declare ENVs in your shell file (`.profile`, `.bashrc` or `.zshrc`)

```sh
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"
```

- Copy `.env.dist` to `.env` and set environments
- `make init`


### Code Quality

Code quality is guarded by pre-commit, but if you should ignore commit checking then you must add the `--no-verify` argument.

### Graphql

How modify gql's schemas && API?:

1. Update schema: api/gql/graph/schemas/
2. Run `make gql-generate` - generate code defined in schemas
3. To update version in config/api-version && optionally creates and push git tag with version:
- `make gql-version-up` - PATCH update, for example 1.0.1 -> 1.0.2

or
- `make MODE=MODE gql-version-up` -
For MODE param choose between MAJOR, MINOR or PATCH - https://semver.org.
