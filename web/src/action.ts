import { ActionFunctionArgs, redirect } from "react-router-dom";
import { API_URL } from "./loader";

export async function loginAction({ request }: ActionFunctionArgs) {
  const form = await request.formData();
  const resp = await fetch(`${API_URL}/api/login`, {
    method: "POST",
    body: form,
  });
  if (resp.status === 401) {
    throw new Response("Unauthorized", { status: 401 });
  }
  if (resp.status !== 200) {
    return null;
  }
  const data = await resp.json();
  const token = data.token;
  localStorage.setItem("token", token);
  return redirect("/");
}

export async function logoutAction() {
  localStorage.removeItem("token");
  return redirect("/");
}

export async function hostAction({ request }: ActionFunctionArgs) {
  const form = await request.formData();
  const body = JSON.stringify(Object.fromEntries(form.entries()));
  const resp = await fetch(`${API_URL}/api/hosts`, {
    method: request.method,
    body: body,
    headers: {
      Authorization: localStorage.getItem("token") || "",
    },
  });
  if (resp.status === 401) {
    throw new Response("Unauthorized", { status: 401 });
  }
  if (!resp.ok) {
    return null;
  }
  return null;
}

export async function cachedAction({ request }: ActionFunctionArgs) {
  const form = await request.formData();
  const body = JSON.stringify(Object.fromEntries(form.entries()));
  const resp = await fetch(`${API_URL}/api/cached`, {
    method: request.method,
    body: body,
    headers: {
      Authorization: localStorage.getItem("token") || "",
    },
  });
  if (resp.status === 401) {
    throw new Response("Unauthorized", { status: 401 });
  }
  if (!resp.ok) {
    return null;
  }
  return null;
}
