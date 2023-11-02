import { Form } from "react-router-dom";

export default function LoginPage() {
  return (
    <div className="p-4 max-w-xl mx-auto pt-24">
      <h1 className="font-bold text-xl">Login</h1>
      <Form className="mt-10 w-full flex flex-col gap-4" method="post">
        <div className="flex items-center gap-2">
          <input type="text" className="w-full" placeholder="Email" />
        </div>
        <div className="flex items-center gap-2">
          <input type="password" className="w-full" placeholder="Password" />
        </div>
        <button
          className="font-semibold bg-zinc-800 text-zinc-200 px-5 py-3"
          type="submit"
        >
          Login
        </button>
      </Form>
      <div className="text-sm hidden">admin</div>
    </div>
  );
}
