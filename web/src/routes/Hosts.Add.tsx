import { Form } from "react-router-dom";

export default function HostsAdd() {
  return (
    <div className="pt-10">
      <Form
        className="mt-10 w-full flex flex-col gap-4"
        method="post"
        action="/hosts"
      >
        <div className="flex items-center gap-2">
          <input
            type="text"
            className="w-full"
            name="host"
            placeholder="Host"
          />
        </div>
        <div className="flex items-center gap-2">
          <input
            type="text"
            className="w-full"
            name="origin"
            placeholder="Origin"
          />
        </div>
        <button
          className="font-semibold bg-zinc-800 text-zinc-200 px-5 py-3"
          type="submit"
        >
          Add
        </button>
      </Form>
    </div>
  );
}
