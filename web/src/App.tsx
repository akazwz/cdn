import { Form, Link, useRouteLoaderData } from "react-router-dom";

export default function App() {
  const rootData = (useRouteLoaderData("root") as { user: string }) || null;
  const authed = rootData && rootData.user;
  return (
    <div className="p-4 max-w-3xl mx-auto">
      <div className="flex items-center">
        <Link className="text-xl font-bold" to="/">
          Dash
        </Link>
        <div className="ml-auto">
          {authed ? (
            <Form action="/logout" method="post" className="flex items-center">
              <button
                type="submit"
                className="bg-zinc-100 px-3 py-2 rounded-md text-red-600 font-semibold"
              >
                Logout
              </button>
            </Form>
          ) : (
            <Link
              to="/login"
              className="bg-zinc-100 px-3 py-2 rounded-md font-semibold text-blue-600"
            >
              Login
            </Link>
          )}
        </div>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-10">
        <Link
          to="/hosts"
          className="p-4 bg-zinc-100 rounded-md font-semibold text-blue-600"
        >
          Hosts
        </Link>
        <Link
          to="/cached"
          className="p-4 bg-zinc-100 rounded-md font-semibold text-blue-600"
        >
          Cached
        </Link>
      </div>
    </div>
  );
}
