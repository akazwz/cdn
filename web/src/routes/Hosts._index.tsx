import { InboxIcon, TrashIcon } from "@heroicons/react/24/outline";
import { Form, useRouteLoaderData } from "react-router-dom";
import { HostOrigin } from "../loader";

export default function HostIndex() {
  const data = useRouteLoaderData("hosts") as HostOrigin[];
  return (
    <div className="pt-10">
      <div className="flex flex-col gap-4">
        {data.map((item) => {
          return (
            <div key={item.host} className="p-4 bg-zinc-100 rounded-md flex">
              <div className="flex flex-col gap-2">
                <div className="font-semibold">{item.host}</div>
                <div className="text-sm">{item.origin}</div>
              </div>
              <div className="flex ml-auto items-center">
                <Form
                  method="delete"
                  action="/hosts"
                  className="flex items-center"
                >
                  <button type="submit">
                    <input
                      type="text"
                      readOnly
                      className="hidden"
                      name="host"
                      value={item.host}
                    />
                    <TrashIcon className="h6 w-6 text-red-600" />
                  </button>
                </Form>
              </div>
            </div>
          );
        })}
        {data.length === 0 && (
          <div className="flex flex-col items-center justify-center pt-10 text-xl font-semibold gap-5">
            <div className="text-gray-400">No hosts found</div>
            <InboxIcon className=" h-24 w-24 text-gray-400" />
          </div>
        )}
      </div>
    </div>
  );
}
