import { Input } from "@/components/Input/Input";
import DefaultLayout from "@/layouts/DefaultLayout";
import { postSignUp } from "@/server/postSignUp";
import { useRouter } from "next/router";
import { useState } from "react";

export default function RegisterPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [password2, setPassword2] = useState("");
  const [email, setEmail] = useState("");
  const router = useRouter();

  return (
    <DefaultLayout>
      <div className='mt-10 bg-gray-300 w-96 rounded-md mx-auto'>
        <form
          className='flex flex-col m-2 p-2'
          onSubmit={(e) => handleSignUp(e)}>
          <Input
            type={"text"}
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            label={"Username"}
          />
          <Input
            type={"email"}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            label={"Email"}
          />
          <Input
            type={"password"}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            label={"Password"}
          />
          <Input
            type={"password"}
            value={password2}
            onChange={(e) => setPassword2(e.target.value)}
            label={"Repeat password"}
          />

          <button
            className='border bg-white rounded-md shadow-md my-4'
            type='submit'>
            Login
          </button>
        </form>
      </div>
    </DefaultLayout>
  );

  async function handleSignUp(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    if (password !== password2) {
      alert("Passwords do not match");
      return;
    }

    const res = await postSignUp({ username, password, password2, email });

    if (res.message) {
      alert(res.message);
      router.push("/");
    } else {
      alert("Something went wrong");
    }
  }
}
