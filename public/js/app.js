// game
const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');

// Set canvas size
canvas.width = 800;
canvas.height = 600;

// components
let players = {};
let myId = null;

// websocket
const socket = new WebSocket(`ws://${window.location.host}/ws`);

socket.onopen = function(event) {
    console.log("WebSocket connection opened");
};

socket.onmessage = function(event) {
    const message = JSON.parse(event.data);
    switch (message.type) {
        case "your_id":
            myId = message.payload;
            break;

        case "update":
            players = message.payload;
            break;

        case "player_left":
            delete players[message.payload];
            break;
    }
};

socket.onerror = function(error) {
    console.error("WebSocket error:", error);
};

socket.onclose = function(event) {
    console.log("WebSocket connection closed:", event);
    alert("server closed");
};

// fire event
document.addEventListener("keydown", (event) => {
    let direction = "";
    
    if (event.key === 'w') {
        direction = 'w';
    } 
    if (event.key === 'a') {
        direction = 'a';
    } 
    if (event.key === 's') {
        direction = 's';
    } 
    if (event.key === 'd') {
        direction = 'd';
    } 

    if (direction) {
        socket.send(JSON.stringify({type: 'movement', payload: direction}));
    }
});

// draw function
function draw() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    for (let id in players) {
        const player = players[id];
        ctx.fillStyle = (id === myId) ? 'blue' : 'red';
        ctx.fillRect(player.x, player.y, 50, 50); // Draw a square for each player
        ctx.fillStyle = 'black';
        ctx.fillText(player.name, player.x, player.y - 10); // Draw player name
    }
}

// game loop
function gameLoop() {
    draw();
    requestAnimationFrame(gameLoop);
}

// start the game loop
requestAnimationFrame(gameLoop);
