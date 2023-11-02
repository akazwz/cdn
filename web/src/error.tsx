import { Link, isRouteErrorResponse, useRouteError } from "react-router-dom";

export function RootBoundary() {
  const error = useRouteError();

  let DefaultError = <div>Something went wrong</div>;

  if (isRouteErrorResponse(error)) {
    switch (error.status) {
      case 404:
        DefaultError = <div>This page doesn't exist!</div>;
        break;
      case 401:
        DefaultError = (
          <div className="flex flex-col items-center gap-4">
            You aren't authorized to see this
            <Link
              to="/login"
              className="bg-blue-600 px-5 py-3 border rounded-md text-white"
            >
              Login
            </Link>
          </div>
        );
        break;
      case 503:
        DefaultError = <div>Looks like our API is down</div>;
        break;
      case 418:
        DefaultError = <div>🫖</div>;
        break;
    }
  }

  return (
    <div className="font-semibold text-red-600 max-w-3xl mx-auto items-center flex flex-col pt-24 gap-10">
      {DefaultError}
      <div className="w-full p-10 max-w-xl">
        <Link
          className=" bg-zinc-800 text-zinc-200 flex px-5 py-5 w-full items-center justify-center"
          to="/"
        >
          Go Home
        </Link>
      </div>
    </div>
  );
}
