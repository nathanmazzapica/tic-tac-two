#### Idle
Event: Exist -> Action: StartTTLTimer -> Next: Idle
Event: TTLTimerExpires: Action Close -> Next: Closed
Event: PlayerConnects -> Action: StopTimer() -> Next: WaitingForSecond

#### WaitingForSecond
Event: SecondConnects -> Action: ? -> Next: InProgress
Event: SecondDisconnects -> Action: ? -> Next: Idle

#### InProgress
Event: Move -> Action: Validate
    Valid:   Apply, broadcast -> Next: InProgress or Terminal {Winner, Method Win|Draw} 
    Invalid: Reply error -> Next: InProgress
Event: Disconnect -> Action: Begin deadline, broadcast {paused:true, missing:player, deadline} -> Next: InProgress

Event: Reconnect -> Action: Broadcast {paused: false}, close deadline -> Next: InProgress
Event: MissedDeadline -> Action: Broadcast -> Next: Terminal {Winner, Method: Forfeit}
Event: Timeout -> Action: broadcast timeout -> Next: Terminal{Winner=None, Method: Timeout}
Event: Error/InvalidMessage -> Action: send error{code, msg} -> Next: InProgress

#### Terminal
Event: Rematch -> Action: ResetBoard() -> Next: InProgress
Event: Disconnect -> Action: Start close timer -> Next: Terminal
Event: CloseTimerExpires -> Action: CloseLobby()
