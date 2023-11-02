import { PlusIcon } from "@heroicons/react/24/outline";
import { Link, Outlet } from "react-router-dom";

export default function HostsLayout() {
  return (
    <div className="p-4 max-w-3xl mx-auto">
      <div className="flex gap-4 items-center">
        <Link className="font-bold text-xl" to="/hosts">
          Hosts
        </Link>
        <Link
          to="/"
          className="text-sm font-semibold bg-zinc-100 p-2 rounded-md"
        >
          Dash
        </Link>
        <div className="flex ml-auto items-center">
          <Link className="bg-zinc-100 p-1 rounded-md" to="/hosts/add">
            <PlusIcon className="h6 w-6" />
          </Link>
        </div>
      </div>
      <Outlet />
    </div>
  );
}
