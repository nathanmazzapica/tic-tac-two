console.log("JS loaded!");

const LobbyState = {
    Idle: 0,
    WaitingForSecond: 1,
    InProgress: 2,
    Terminal: 3,
    Closed: 4,
}

const GameState = {
    InProgress: 0,
    Won: 1,
    Draw: 2,
}

const GameMark = {
    Empty: 0,
    X: 1,
    O: 2,
}

const ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = () => {
    console.log("Connected to server successfully")
}

ws.onmessage = (e) => {
    console.log("Received message from server")
}

ws.onclose = () => {
    alert("Connection closed")
}

function testCommand() {
    let join = {
        type: 'join',
        data: {
            player_id: '1'
        }
    }

    ws.send(JSON.stringify(join))
}