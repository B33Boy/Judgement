# Reconnect support: stable player identity + disconnect grace period

Status: **planned, not implemented.**

## Context

Right now a dropped WebSocket is indistinguishable from a player quitting:
`NewPlayer` mints a fresh UUID per connection (`internal/app/player.go`),
and `ws.go`'s deferred cleanup calls `RemovePlayer` immediately on any
disconnect — deleting the player from the session and, if they're the last
one, tearing down the whole session. A wifi blip, a backgrounded tab, or a
page refresh permanently ends that player's game.

This implements **server-authoritative state + reconnect tokens with a
grace period** — the pattern casual multiplayer games (Kahoot, Jackbox-style)
use. It does *not* attempt to survive the server process itself restarting
(that needs external state persistence — a much bigger lift, deliberately
out of scope) — only network-level disconnects of an otherwise-alive process.

**Security note driving the design:** the existing `PlayerID` is already
broadcast to everyone via `players_update` (`PlayerPublic{ID, Name}`), so
it cannot double as a reconnect secret — any player could otherwise hijack
another's seat by replaying their known ID. The plan introduces a
separate, never-broadcast `ReconnectToken` sent privately only to its owner.

## Approach

### 1. Backend types (`internal/types/types.go`)
Add `ReconnectToken string` to `Player`. Stays server-internal — never
included in `PlayerPublic` or any broadcast payload.

### 2. `internal/app/player.go`
`NewPlayer` mints both `ID` and `ReconnectToken` (`uuid.NewString()` each)
for a fresh join, as today, just with the extra field.

### 3. `internal/app/session.go` — grace-period bookkeeping
- Add `pendingRemoval map[t.PlayerID]*time.Timer` and a `graceDuration`
  field (default constant `reconnectGrace = 90 * time.Second`, overridable
  by tests via an unexported constructor so tests don't sleep 90s).
- **`HandleDisconnect(player *t.Player)`** replaces the current
  unconditional `RemovePlayer` call from `ws.go`'s defer: if
  `s.players[player.ID] == player` (pointer check — guards the race where
  a reconnect already replaced this entry before the old connection's
  cleanup runs, so a stale disconnect can't clobber a fresh one), start a
  timer that calls today's full removal logic
  (cancel/close/delete/broadcast/maybe-cancel-session) after
  `graceDuration`. Otherwise no-op — a newer connection already took over.
- **`Reconnect(token string, conn *websocket.Conn) (*t.Player, bool)`**:
  scan `s.players` for a matching `ReconnectToken`. Only succeeds if that
  player is currently in `pendingRemoval` (mid-grace-period) — a token for
  someone still actively connected, or whose grace period already expired,
  falls through to a normal fresh join instead. On success: stop/remove
  the timer, build a new `*t.Player` reusing the existing `ID` +
  `ReconnectToken` but fresh `Conn`/`Send`/`Ctx`/`Cancel`, swap it into
  `s.players`, return it.
- This means a reconnecting player transparently bypasses the
  `gameStarted` late-join rejection (`AddPlayer`) — intentional: they're
  not a new joiner, they're resuming a seat that already exists in
  `game.Players`. A token that doesn't match anyone mid-grace-period still
  goes through `AddPlayer` and gets the existing rejection if the game has
  started — that gate is unchanged for genuinely new joiners.

### 4. `internal/game` — resync on reconnect
`Game` methods are only ever safe to call from `Session.run`'s single
goroutine (see `ARCHITECTURE.md`'s concurrency note) — the ws-handling
goroutine can't call into `Game` directly on reconnect. Instead:
- Add a server-internal message type `t.MsgResync` (never sent by the
  frontend — synthesized by `ws.go` after a successful `Reconnect` and
  pushed through the existing `s.inputs` channel, so it's processed by
  the same goroutine that owns `Game`).
- `session.go: handleInput` gets a case for it: `if s.game != nil {
  s.game.Resync(input.Player.ID) }`.
- New method in `internal/game/events.go`: `func (g *Game)
  Resync(id t.PlayerID)` — looks up `g.Players[id]`, calls the existing
  `sendCardsToPlayer` (their hand) and `broadcastGameState` (current
  state; broadcasting to everyone instead of just the reconnecting player
  is redundant for the rest but harmless, and reuses what's already there).

### 5. `internal/app/ws.go`
- Read an optional `reconnectToken` query param alongside
  `sessionId`/`playerName`.
- If present, try `session.Reconnect(token, conn)` first. On success: skip
  `AddPlayer` entirely, send `welcome` (same `ID`/`ReconnectToken` — the
  client already has them, but resending keeps the flow uniform), broadcast
  `players_update`, and push a `t.MsgResync` input for this player.
- Otherwise (no token, or `Reconnect` returned false): today's fresh-join
  path (`NewPlayer` → `AddPlayer`, existing rejection handling unchanged).
- Change the deferred cleanup from `onPlayerLeave` (immediate
  `RemovePlayer`) to `session.HandleDisconnect(player)`.
- `welcome`'s payload changes from a bare `player.ID` string to
  `{playerId, reconnectToken}` (small struct in `internal/app/payload.go`).
- Add `Connected bool` to `PlayerPublic`/`broadcastPlayersUpdate` (derived
  from whether that player is currently in `pendingRemoval`) so the lobby
  can show who's mid-disconnect — cheap, and otherwise a stalled turn on a
  disconnected player looks like nothing is happening.

### 6. Frontend (`frontend/src`)
- `lib/player.ts`: add `getReconnectToken(sessionId)` /
  `setReconnectToken(sessionId, token)`, keyed by session ID in
  `localStorage` (`reconnectToken:<sessionId>`) so tokens from different
  games don't collide.
- `context/GameContext.tsx`:
  - `connect(sessionId, playerName)` appends `&reconnectToken=...` to the
    WS URL when one is stored for that session.
  - `"welcome"` handler reads `{playerId, reconnectToken}`, persists the
    token via `setReconnectToken`.
  - `ws.onclose`: track whether the close was from an explicit
    `disconnect()` call (a ref flag). If not, retry `connect(...)` on a
    timer — every 3s, up to 30 attempts (~90s total, matching
    `reconnectGrace` so the client keeps trying for exactly as long as the
    server will hold the seat). Clear the retry loop on a successful
    reopen or on unmount.
- `types.ts`: update `WSEnvelope`'s `welcome` payload shape; add
  `connected: boolean` to `PlayerPublic`.
- `pages/Session.tsx`: render disconnected players dimmed / with a
  "(reconnecting…)" suffix using the new `connected` field.

## Explicitly out of scope
- Surviving the server process restarting (needs external state
  persistence — separate, bigger effort, deliberately deferred).
- Auto-skipping a disconnected player's turn — the game just waits during
  their grace period, same as today waiting on any slow player. A
  turn-timeout feature is a separate concern.

## Verification (once implemented)
1. `go build ./...`, `go vet ./...`, `go test ./... -race`.
2. New tests in `internal/app/session_test.go`: disconnect starts a grace
   timer; reconnect within the grace period reattaches to the same
   `PlayerID` (and the pointer-identity race guard); reconnect after grace
   expiry does not reattach and falls through to normal join rules; a
   token presented while the original connection is still active is
   rejected. Use a short injectable grace duration so these don't sleep
   90s for real.
3. New test in `internal/game` for `Game.Resync` (hand + state resent to
   the right player).
4. Live check with a browser-driving harness (Playwright, as used earlier
   in this project): start a 3-player game, forcibly close one player's
   WebSocket mid-game (not via the UI's disconnect button), reconnect that
   same browser context with the stored token, and confirm their
   hand/game state comes back and they can keep playing.
