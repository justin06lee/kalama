# kalama

The engine library behind the [shaw](https://github.com/justin06lee/shaw) terminal
arcade. Games import `github.com/justin06lee/kalama` and implement the `Game`
interface; kalama owns the terminal, runs a fixed-timestep loop, and renders a
truecolor pixel canvas using half-block glyphs.

Players never type `kalama` — they use the [`shaw`](https://github.com/justin06lee/shaw)
launcher. This module is for game developers.

## Use it

```go
import "github.com/justin06lee/kalama"
```

Implement `kalama.Game` (`Update(dt, in) Action` + `Draw(*Canvas)`) and call
`kalama.Run(game, kalama.Options{})`.

## Persistence

`kalama.DataDir("<game>.shaw")` returns a per-game directory under
`~/.shaw/data/` (or `$SHAW_DATA_DIR`).

## Develop

```
make test
make vet
make tidy
```

## Related

- [shaw](https://github.com/justin06lee/shaw) — the launcher/package manager players run.
- [hegale](https://github.com/justin06lee/hegale) — the game registry.
