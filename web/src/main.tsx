import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { cachedLoader, hostLoader, rootLoader } from "./loader";
import Cached from "./routes/Cached";

import { cachedAction, hostAction, loginAction, logoutAction } from "./action";
import App from "./App";
import { RootBoundary } from "./error";
import "./index.css";
import HostsLayout from "./routes/Hosts";
import HostIndex from "./routes/Hosts._index";
import HostsAdd from "./routes/Hosts.Add";
import Login from "./routes/Login";
import Me from "./routes/Me";

const router = createBrowserRouter([
  {
    path: "/",
    id: "root",
    errorElement: <RootBoundary />,
    loader: rootLoader,
    children: [
      {
        path: "",
        element: <App />,
      },
      {
        path: "me",
        id: "me",
        element: <Me />,
      },
      {
        path: "login",
        element: <Login />,
        action: loginAction,
      },
      {
        path: "hosts",
        id: "hosts",
        element: <HostsLayout />,
        loader: hostLoader,
        action: hostAction,
        children: [
          {
            path: "",
            element: <HostIndex />,
          },
          {
            path: "add",
            element: <HostsAdd />,
          },
        ],
      },
      {
        path: "/cached",
        element: <Cached />,
        loader: cachedLoader,
        action: cachedAction,
      },
    ],
  },
  {
    path: "/logout",
    action: logoutAction,
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
