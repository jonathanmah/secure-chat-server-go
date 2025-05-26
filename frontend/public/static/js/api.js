import { SERVER_BASE_URL } from "./config.js";

// GET JSON - the current authenticated users id and username
export async function getUserInfo() {
  const res = await fetch(`${SERVER_BASE_URL}/auth/user-info`, {
    credentials: "include", // include cookies
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(
      `Failed to fetch user (${res.status} ${res.statusText}): ${errorText}`
    );
  }
  const data = await res.json();
  return data;
}

// POST username - update the current users username in db
export async function updateUsername(newUsername) {
  const res = await fetch(`${SERVER_BASE_URL}/auth/update-username`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: newUsername }),
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(
      `Failed to update username (${res.status} ${res.statusText}): ${errorText}`
    );
  }
}

// POST empty - invalidate session cookie and redirects to login
export async function logout() {
  const res = await fetch(`${SERVER_BASE_URL}/auth/logout`, {
    method: "POST",
    credentials: "include",
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(
      `Failed to logout (${res.status} ${res.statusText}): ${errorText}`
    );
  }
}

// POST email - send link with token for password reset
export async function sendPasswordResetEmail(email) {
  const res = await fetch(`${SERVER_BASE_URL}/auth/forgot-password`, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      email: email,
    }),
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(
      `Failed to send reset link (${res.status} ${res.statusText}): ${errorText}`
    );
  }
}

// POST email, password to login
export async function login(email, password) {
  const res = await fetch(`${SERVER_BASE_URL}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    credentials: "include", // add cookie
    body: new URLSearchParams({
      email,
      password,
    }),
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(
      `Failed to login(${res.status} ${res.statusText}): ${errorText}`
    );
  }
}

// POST token, password to reset
export async function resetPassword(token, password) {
  const res = await fetch(`${SERVER_BASE_URL}/auth/reset-password`, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      token: token,
      password: password,
    }),
  });
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(
      `Failed to reset password(${res.status} ${res.statusText}): ${errorText}`
    );
  }
}

// POST email, password, register a new account
export async function createNewAccount(email, password) {
  const res = await fetch(`${SERVER_BASE_URL}/auth/sign-up`, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      email: email,
      password: password,
    }),
  });
  console.log("AFTER FETCH");
  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(
      `Failed to create account(${res.status} ${res.statusText}): ${errorText}`
    );
  }
}
