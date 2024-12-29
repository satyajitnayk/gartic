// DOM Elements
const roomInput = document.getElementById('room-name');
const userInput = document.getElementById('user-name');
const connectBtn = document.getElementById('connect-btn');
const sendMessageBtn = document.getElementById('send-msg-btn');
const formContainer = document.getElementById('form-container');
const gameContainer = document.getElementById('game-container');
const guessInput = document.getElementById('guess-input');
const messages = document.getElementById('messages');

const roomDisplay = document.getElementById('room-name-display');
const userDisplay = document.getElementById('user-name-display');

const canvas = document.getElementById('drawing-board');
const ctx = canvas.getContext('2d');

// WebSocket Connection
let ws;
let drawing = false;

// Event Listeners
connectBtn.addEventListener('click', handleConnect);
canvas.addEventListener('mousedown', () => (drawing = true));
canvas.addEventListener('mouseup', () => (drawing = false));
canvas.addEventListener('mousemove', handleDrawing);
guessInput.addEventListener('keypress', handleGuess);
sendMessageBtn.addEventListener('click', sendGuess);

// Update canvas size on window resize (if needed)
window.addEventListener('resize', setCanvasSize);

// Functions

function setCanvasSize() {
  const containerWidth = gameContainer.offsetWidth;
  const containerHeight = gameContainer.offsetHeight;

  // Set the actual drawing canvas size to match the container size
  canvas.width = containerWidth;
  canvas.height = containerHeight * 0.7;

  // Set the CSS size to match the drawing size (this ensures no scaling)
  canvas.style.width = `${containerWidth}px`;
  canvas.style.height = `${containerHeight}px`;
}

/**
 * Handle WebSocket connection and UI updates when the connect button is clicked
 */
function handleConnect() {
  const roomName = roomInput.value.trim();
  const userName = userInput.value.trim();

  if (!roomName || !userName) {
    alert('Please enter both room name and username.');
    return;
  }

  connectToWebSocket(roomName, userName);

  // WebSocket Event Handlers
  ws.onopen = () => {
    console.log('WebSocket connected');
    ws.send(JSON.stringify({ type: 'join', room: roomName, name: userName }));
  };

  ws.onmessage = handleWebSocketMessage;
  ws.onerror = (error) => console.error('WebSocket error:', error);

  roomDisplay.textContent = roomName;
  userDisplay.textContent = userName;

  // Update UI
  formContainer.style.display = 'none';
  gameContainer.style.display = 'flex';

  // Resize the canvas after the room is joined
  setTimeout(setCanvasSize, 0); // Delay the resizing for DOM update
}

function connectToWebSocket(room, username) {
  ws = new WebSocket(`ws://localhost:8080/ws?room=${room}&name=${username}`);
}

/**
 * Handle incoming WebSocket messages
 * @param {MessageEvent} event - WebSocket message event
 */
function handleWebSocketMessage(event) {
  const msg = JSON.parse(event.data);

  if (msg.type === 'draw') {
    const { x, y } = JSON.parse(msg.content);
    drawOnCanvas(x, y, false); // Draw from received data
  } else if (msg.type === 'guess') {
    displayGuess(msg.content);
  }
}

/**
 * Handle drawing on the canvas
 * @param {MouseEvent} e - Mouse move event
 */
function handleDrawing(e) {
  if (!drawing) return;

  const rect = canvas.getBoundingClientRect(); // Get the canvas position relative to the viewport
  // Calculate the mouse position relative to the canvas
  const x = (e.clientX - rect.left) * (canvas.width / rect.width); // Adjust for scaling
  const y = (e.clientY - rect.top) * (canvas.height / rect.height); // Adjust for scaling

  // Draw locally and broadcast
  drawOnCanvas(x, y, true);
}

/**
 * Draw on the canvas
 * @param {number} x - X-coordinate
 * @param {number} y - Y-coordinate
 * @param {boolean} broadcast - Whether to send drawing data via WebSocket
 */
function drawOnCanvas(x, y, broadcast = true) {
  ctx.lineTo(x, y);
  ctx.stroke();
  ctx.beginPath();
  ctx.moveTo(x, y);

  if (broadcast) {
    const message = { type: 'draw', content: JSON.stringify({ x, y }) };
    ws.send(JSON.stringify(message));
  }
}

function sendGuess() {
  const guess = guessInput.value.trim();
  if (guess) {
    const message = { type: 'guess', content: guess };
    ws.send(JSON.stringify(message));
    guessInput.value = '';
  }
}

function handleGuess(e) {
  if (e.key === 'Enter') {
    e.preventDefault(); // Prevent form submission or other default actions
    sendGuess();
  }
}

/**
 * Display a guess in the message list
 * @param {string} guess - The guess to display
 */
function displayGuess(guess) {
  const li = document.createElement('li');
  li.textContent = guess;
  messages.appendChild(li);
}

// Make sure the canvas is properly resized when the window is loaded
window.onload = setCanvasSize;
