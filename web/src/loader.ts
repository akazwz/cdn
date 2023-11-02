export type HostOrigin = {
  host: string;
  origin: string;
};

export const API_URL = import.meta.env.VITE_API_URL as string;

export async function rootLoader() {
  const resp = await fetch(`${API_URL}/api/me`, {
    headers: {
      Authorization: localStorage.getItem("token") || "",
    },
  });
  if (!resp.ok) {
    return null;
  }
  const data = await resp.json();
  return data;
}

export async function hostLoader() {
  const resp = await fetch(`${API_URL}/api/hosts`, {
    headers: {
      Authorization: localStorage.getItem("token") || "",
    },
  });
  if (resp.status === 401) {
    throw new Response("Unauthorized", { status: 401 });
  }
  const data = await resp.json();
  return data;
}

export type CachedType = {
  cache_key: string;
  cache_path: string;
};

export async function cachedLoader() {
  const resp = await fetch(`${API_URL}/api/cached`, {
    headers: {
      Authorization: localStorage.getItem("token") || "",
    },
  });
  if (resp.status === 401) {
    throw new Response("Unauthorized", { status: 401 });
  }
  const data = await resp.json();
  return data;
}
