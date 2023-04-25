import { useForm } from "@/hooks/useForm";
import { useLogin } from "@/hooks/useLogin";
import useUserStore from "@/store/userStore";
import { LoginSchema } from "@/types/Auth.type";
import { useState } from "react";
import { useCookies } from "react-cookie";
import { useRouter } from "next/router";
import { Button } from "../Input/Button";

export function Login() {
	const [, setUserIdCookie] = useCookies(["user_id"]);
	const setUserIdStore = useUserStore((state) => state.setUserId);
	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");

	const router = useRouter();

	const { mutate: mutateLogin, isLoading } = useLogin();

	const loginForm = useForm(LoginSchema, (values) => {
		mutateLogin(values, {
			onSuccess: (mutationResp) => {
				if (mutationResp?.user_id) {
					setUserIdCookie("user_id", mutationResp?.user_id);
					setUserIdStore(mutationResp?.user_id);
				}
				router.push("/");
			},
			onError: (error) => {
				alert(error);
			},
		});
	});
	return (
		<form
			className="flex flex-col p-4 dark:bg-slate-900"
			onSubmit={(e) => {
				e.preventDefault();
				loginForm.onSubmit({
					username,
					password,
				});
			}}
		>
			<input
				className="px-2 py-1 my-4 shadow-md rounded-md"
				placeholder="Username"
				type="text"
				value={username}
				onChange={(e) => setUsername(e.target.value)}
			/>
			<input
				className="px-2 py-1 my-4 shadow-md rounded-md"
				placeholder="Password"
				type="password"
				value={password}
				onChange={(e) => setPassword(e.target.value)}
			/>
			<Button isLoading={isLoading} text={"Login"} />
		</form>
	);
}
