import {
  roomHeader,
  chatMessages,
  messageInput,
  usernameDisplay,
  darkModeToggle,
} from "./dom.js";

// -------------------------------------- CHAT MESSAGE DISPLAY ----------------------------------
export function renderRoomHeader(roomID) {
  if (roomHeader) {
    roomHeader.textContent = `Hub - Room ${roomID}`;
  }
}

// render the chat messages
export function renderChatMessage(payload) {
  const messageDiv = createChatMessage(payload);
  chatMessages.append(messageDiv);
  chatMessages.scrollTop = chatMessages.scrollHeight; // scroll to bottom
}

// clear previous messages
export function clearChatMessages() {
  chatMessages.textContent = "";
}

// creates HTML element for chat message
function createChatMessage(payload) {
  const messageDiv = document.createElement("div");
  const timestampSpan = document.createElement("span");
  timestampSpan.classList.add("timestamp");
  const date = new Date(payload.time);
  const time = date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
  timestampSpan.textContent = time;

  if (payload.sender_id === "notification") {
    messageDiv.classList.add("chat-notification");
    messageDiv.textContent = payload.text;
    messageDiv.appendChild(timestampSpan);
  } else {
    messageDiv.classList.add("chat-message");
    const usernameStrong = document.createElement("strong");
    usernameStrong.style.color = getUserColour(payload.sender_username);
    usernameStrong.textContent = payload.sender_username;
    const textNode = document.createTextNode(`: ${payload.text} `);
    messageDiv.append(usernameStrong, textNode, timestampSpan);
  }
  return messageDiv;
}

// -------------------------------------- USER LIST ----------------------------------
// update the user list with currently active users
export function renderActiveUsers(users) {
  const tbody = document.querySelector("#activeUsersTable tbody");
  tbody.textContent = "";
  users.forEach((user) => {
    const tr = document.createElement("tr");
    const td = document.createElement("td");
    const dot = document.createElement("span");
    dot.classList.add("online-dot");
    td.appendChild(dot);
    td.append(` ${user.username}`);
    tr.appendChild(td);
    tbody.appendChild(tr);
  });
}

// -------------------------------------- CHAT MESSAGE INPUT ----------------------------------
const MAX_CHARS = 500;
// updates character limit
export function renderCharCount() {
  charCount.textContent = `${messageInput.value.length} / ${MAX_CHARS}`;
}
// resizes input box
export function resizeTextarea() {
  messageInput.style.height = "auto";
  messageInput.style.height = messageInput.scrollHeight + "px";
}

export function renderUsername(username) {
  usernameDisplay.textContent = username;
}

// -------------------------------------- COLOURS ---------------------------------------------
export function loadDarkModePref() {
  const darkPref = localStorage.getItem("darkMode") === "true";
  darkModeToggle.checked = darkPref;
  document.body.classList.toggle("dark-mode", darkPref);
}

const colours = [
  "#c31442", // red
  "#339d3f", // green
  "#e6cb16", // yellow
  "#3a56bb", // blue
  "#d9722b", // orange
  "#7f1aa0", // purple
  "#3fc7c7", // cyan
  "#c92bc5", // pink
  "#a5d509", // lime
  "#d1a3a3", // light pink
  "#006666", // teal
  "#cda7e6", // lavender
  "#80531f", // brown
  "#d9d3ac", // light lellow
  "#660000", // maroon
];

const userColours = {};

function getUserColour(username) {
  if (!userColours[username]) {
    const availableColors = colours.filter(
      (c) => !Object.values(userColours).includes(c)
    );
    const palette = availableColors.length > 0 ? availableColors : colours;
    const color = palette[Math.floor(Math.random() * palette.length)];
    userColours[username] = color;
  }
  return userColours[username];
}
