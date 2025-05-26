import { createNewAccount } from "../api.js";
import {
  confirmPasswordElement,
  emailElement,
  passwordElement,
} from "./dom.js";

document
  .getElementById("sign-up-form")
  .addEventListener("submit", async function (e) {
    e.preventDefault();
    const email = emailElement.value;
    const password = passwordElement.value;
    const confirm = confirmPasswordElement.value;
    const error = document.getElementById("error");
    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
    if (password !== confirm) {
      error.textContent = "Passwords do not match.";
      return;
    }
    if (!passwordRegex.test(password)) {
      error.textContent =
        "Password must be at least 8 characters long and include uppercase, lowercase, number, and special character.";
      return;
    }
    error.textContent = "";
    try {
      await createNewAccount(email, password);
    } catch (err) {
      error.textContent = err.message;
      return;
    }
    const form = document.getElementById("sign-up-form");
    form.style.display = "none";
    const successMsg = document.createElement("p");
    successMsg.textContent =
      "Account created! Please check your email to confirm your account.";
    successMsg.style.color = "green";
    document.querySelector(".container").appendChild(successMsg);
  });
