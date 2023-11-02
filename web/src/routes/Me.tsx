import { Link, useLoaderData } from "react-router-dom";

export default function Me() {
  const data = useLoaderData() as { user: string };
  return (
    <div className="p-4 max-w-3xl mx-auto">
      <div className="flex gap-4 items-center">
        <Link className="font-bold text-xl" to="/cached">
          Me
        </Link>
        <Link
          to="/"
          className="text-sm font-semibold bg-zinc-100 p-2 rounded-md"
        >
          Dash
        </Link>
      </div>
      <div className="flex items-center justify-center text-xl font-bold pt-24 gap-4">
        You are <span className="bg-blue-100 rounded-md p-2">{data.user}</span>
      </div>
    </div>
  );
}
