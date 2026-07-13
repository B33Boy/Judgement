# Architecture

Judgement (Oh Hell / Diminishing Whist) implemented as a Go WebSocket backend
driving a React/TypeScript frontend. This doc is a map of the system, not a
tutorial — read it alongside the code it points to.

## Process wiring

```
cmd/api/main.go
  -> internal/server.NewServer(app)   http.Server config (port, timeouts)
  -> internal/app.NewApp()            builds the SessionStore, registers routes
       -> RegisterRoutes()            chi router: /api/session, /api/session/{id}, /ws
```

- `POST /api/session` — creates a `Session`, returns its ID.
- `GET /api/session/{id}` — 200/404, used by the frontend to validate a
  session before joining.
- `GET /ws?sessionId=...&playerName=...` — the only real traffic path;
  everything after the handshake goes over this one WebSocket.

## Session layer (`internal/app`)

A `Session` (`session.go`) owns one `*game.Game` and one goroutine
(`Session.run`) that is the single writer to all shared game state — the
game engine itself is **not** goroutine-safe by design, so this loop is
what makes it safe. It selects on three channels:

- `s.inputs` — `t.GameInput` from any player's read loop (`ws.go`), routed
  to `game.HandleGameInput`.
- `s.outputs` — `t.GameOutput` the game engine emits, fanned out to each
  named recipient's `player.Send` channel.
- `s.ctx.Done()` — session teardown (cancelled once the last player leaves).

Each connected player (`ws.go`) gets two goroutines of its own: a read loop
that decodes incoming `t.Envelope`s into `s.inputs`, and a write loop that
drains `player.Send` back out over the socket. `player.ID` is a fresh UUID
minted per connection (`player.go`) — there's no persistent player identity
across reconnects.

`SessionStore` (`session_store.go`) is just a mutex-guarded
`map[string]*Session` keyed by an 8-character random ID.

## Game engine (`internal/game`)

`Game` (`game.go`) is constructed once per session, when the first
`start_game` message arrives (`session.go: handleInput`). It owns:

- `sm *StateMachine` — a small generic FSM (`state_machine.go`): states +
  events -> next state, nothing game-specific.
- `cycler *PlayerCycler` (`cycler.go`) — turn order. Used for **both**
  bidding order and play order within a trick; `StartFrom(id)` re-anchors
  it, `Next()` advances, `CompletedCycle()`/`WillCompleteNext()` detect the
  first/last player in the current lap.
- `state *GameState` — the piece that gets JSON-marshalled and broadcast
  to every client on every change (`events.go: broadcastGameState`).
- `scores PlayerScore` — `map[PlayerID][]Score`, one slot per round,
  aliased into `state.Scores` so it's included in every broadcast for free.

### State machine

```
StateBid --(BiddingDone)--> StatePlay --(PlayingDone)--> StateResolution
                                ^                              |
                                |----------(TrickContinue)-----|  (more tricks left this round)
                                                                |
StateBid <----(PlayingContinue, i.e. round continues)----------|
                                                                |
StateGameOver <-----------------------(GameDone)----------------
```

`state_transition.go: onStateChanged` runs side effects on entry to each
state — critically, entering `StateResolution` calls `g.handleResolution()`
**synchronously**, right there, because resolution has no client-driven
input to wait for (there's no "ack the trick" message from the frontend).

### Round lifecycle (`round.go`, `handlers.go`)

- **Deal** (`startNewRound`): hand size comes from the pyramid
  `roundSchedule` in `game.go` — `7,6,5,4,3,2,1,1,2,3,4,5,6,7` (14 rounds).
  The round starter rotates one seat (`g.dealerIdx`) each round; trump,
  bids, hands-won, and the table are all reset to zero here.
- **Bid** (`handleBid` / `validateBid`): bid must be `0..cardsInHand`, and
  if `cycler.WillCompleteNext()` (this player is last to bid), the bid
  can't make the total of all bids equal the cards dealt — the classic
  "hook" rule. Rejections go out as `invalid_action` via
  `sendInvalidMove`, and the player's turn does not advance.
- **Trump**: the suit of the very first card led each round becomes trump
  (`handleNoTrumpSuit`, idempotent once set) — there's no flipped card or
  fixed rotation.
- **Play** (`handlePlay` / `isCardPlayable`): must follow the suit that
  was **led** (`g.cardstack[0].Suit`, not the most recently played card)
  if the player holds one; otherwise anything is legal, trump included.
- **Resolve** (`handleResolution` / `determineTrickWinner` / `beats`):
  trump beats non-trump; among two trumps or two led-suit cards, higher
  rank wins; anything else can't win. Winner leads the next trick.
- **Score** (`scoreRound`, called from `finishRound` before `state.Round`
  is incremented): `10 + bid` if a player's bid exactly matches tricks
  won that round, else `0`.
- Game ends when `state.Round` reaches `len(roundSchedule)`.

### Concurrency note

Everything under `internal/game` assumes it's only ever called from
`Session.run`'s single goroutine. There's no locking inside `Game` itself
— don't call `Game` methods from anywhere else without going through the
session's channels.

## Wire protocol (`internal/types/types.go`)

Every message is `{"type": MessageType, "payload": json.RawMessage}`.

| Direction | Type | Payload | Notes |
|---|---|---|---|
| FE -> BE | `start_game` | — | first player-independent action; builds the `Game` |
| FE -> BE | `make_bid` | `{bid}` | only honored during `StateBid`, from `state.turnPlayer` |
| FE -> BE | `play_card` | `{suit, rank}` | only honored during `StatePlay` |
| BE -> FE | `welcome` | `PlayerID` | sent once, right after a socket connects |
| BE -> FE | `players_update` | `PlayerPublic[]` | broadcast on join/leave |
| BE -> FE | `game_started` | — | tells the frontend to route to `/game/:sessionId/:playerName` |
| BE -> FE | `player_hand` | `{cards: string[]}` | private, sent only to that player |
| BE -> FE | `state_sync` | `GameState` (incl. `scores`) | broadcast after every mutation |
| BE -> FE | `invalid_action` | `{message}` | rejected bid/play, sent only to the offending player |
| BE -> FE | `game_end` | — | sent once, on entering `StateGameOver` |

`GameState` and `types.ts` are hand-kept in sync — there's no shared schema
or codegen, so a backend field change means a manual `frontend/src/types.ts`
edit.

## Frontend (`frontend/src`)

- `context/GameContext.tsx` — owns the one WebSocket, dispatches incoming
  envelopes into React state (`gameState`, `hand`, `players`, `playerId`).
  Everything else reads from this context via `useGame()`.
- `pages/Game.tsx` — assembles the game screen from `gameState.state`;
  drag-and-drop card play goes through `@dnd-kit` (`handleDragEnd` ->
  `sendMessage("play_card", ...)`).
- `components/` — `BidBox`, `GameTable`, `PlayerHand`, `ScoreTable`,
  `RoundData`, `SessionBox`, each a thin view over `gameState`/`hand`.

## Known gaps (as of this writing)

- `BidBox` always renders bid buttons `0..7` regardless of the actual
  hand size for the round — backend now rejects out-of-range bids, so
  it's safe but confusing UX on small rounds.
- `GameTable.tsx` has a real type error: `gameState.table[p.id]` is typed
  `string | undefined` but `TableEntry` expects a `Card` — pre-existing,
  not yet fixed.
- A handful of pre-existing lint errors (`no-explicit-any` on the message
  payload types, a conditionally-called `useEffect` in `Game.tsx`).
- No reconnect/resume story — a dropped socket loses that player's seat
  (a fresh UUID is minted on reconnect, so the game engine has no way to
  recognize "the same person came back").
