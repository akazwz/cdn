import { InboxIcon, TrashIcon } from "@heroicons/react/24/outline";
import { Form, Link, useLoaderData } from "react-router-dom";
import { CachedType } from "../loader";

export default function CachedPage() {
  const data = useLoaderData() as CachedType[];
  return (
    <div className="p-4 max-w-3xl mx-auto">
      <div className="flex gap-4 items-center">
        <Link className="font-bold text-xl" to="/cached">
          Cached
        </Link>
        <Link
          to="/"
          className="text-sm font-semibold bg-zinc-100 p-2 rounded-md"
        >
          Dash
        </Link>
      </div>
      <div className="mt-10 flex flex-col gap-2">
        {data.map((item) => (
          <div
            key={item.cache_key}
            className="bg-zinc-100 p-2 rounded-md flex items-center gap-2"
          >
            <p className="font-semibold text-sm">{item.cache_key}</p>
            <div className="flex ml-auto items-center">
              <Form method="delete" className="flex items-center">
                <input
                  type="text"
                  readOnly
                  className="hidden"
                  name="cache_key"
                  value={item.cache_key}
                />
                <button>
                  <TrashIcon className="w-6 h-6 text-red-600" />
                </button>
              </Form>
            </div>
          </div>
        ))}
        {data.length === 0 && (
          <div className="flex flex-col items-center justify-center pt-10 text-xl font-semibold gap-5">
            <div className="text-gray-400">No cached found</div>
            <InboxIcon className=" h-24 w-24 text-gray-400" />
          </div>
        )}
      </div>
    </div>
  );
}
