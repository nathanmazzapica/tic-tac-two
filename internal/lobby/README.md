### States
`Idle > Waiting > InProgress > Terminal > Closed`

### Outcome
`{ Winner: X|O|None, Method: Win|Draw|Forfeit|Timeout, Line? }`

### Buffers & deadlines
`broadcast=64`
`send=32` (if full; disconnect)
`writeDeadline=5s`
`pongTimeout=20s`

### Join Rules
first join goes slot0

second goes slot1

3+ gets rejected

### Turn Rules
Slot0=X starts; server authoritative; clients send expectedMove

### Disconnect policy
On drop set `slot.connected=false`, `slot.reconnectDeadline=now+GracePeriod`. If deadline expires; `Terminal{Forfeit}`

If disconnect during WaitingForSecond `WaitingForSecond->Idle`. Start TTL timer for `Lobby`

### Close Policy
`Terminal` > close lobby after 30s; both disconnected in `InProgress` > close after 2m idle TTL

### Events to emit
`Joined{slot, mark}`

`State{board, turn, moveNum, outcome?}`

`Paused{missing, deadline}`

`Resumed{}`

`Error{code,msg}`

### Considerations for later

- Timeout player after `timeLimit` and -> `Terminal{Forfeit}`
- Do slots need to be explicitly cleared if a player leaves during `WaitingForSecond` or `Terminal`?
- Do we need to broadcast anything on `handleLeave()` when `l.state=WaitingForSecond`? Nobody would hear it anyway...

