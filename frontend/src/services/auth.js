import api from "./api";

export async function login(email, password) {
  const response = await api.post("/auth/login", {
    email,
    password,
  });

  const {
    access_token: accessToken,
    token,
    user,
  } = response.data;
  const authToken = accessToken || token;

  if (!authToken || !user) {
    throw new Error("Invalid login response from server");
  }

  localStorage.setItem("token", authToken);
  localStorage.setItem("user", JSON.stringify(user));

  return user;
}

export function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
}

export function getCurrentUser() {
  const user = localStorage.getItem("user");
  return user ? JSON.parse(user) : null;
}