# shaw

A terminal arcade engine. Games import `github.com/justin06lee/shaw` and
implement the `Game` interface; shaw owns the terminal, runs a fixed-timestep
loop, and renders a truecolor pixel canvas using half-block glyphs.

See `docs/superpowers/specs/2026-05-23-arcade-engine-design.md` for the design.

## Development

```bash
make test   # go test ./...
make vet    # go vet ./...
make tidy   # go mod tidy
```
