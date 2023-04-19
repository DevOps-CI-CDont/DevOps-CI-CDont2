import { useState } from "react";
import { Input } from "../Input/Input";
import { useRegister } from "@/hooks/useRegister";
import { useForm } from "@/hooks/useForm";
import { RegisterSchema } from "@/types/Auth.type";
import { useLogin } from "@/hooks/useLogin";
import { useRouter } from "next/router";
import { useCookies } from "react-cookie";
import useUserStore from "@/store/userStore";
import { Button } from "../Input/Button";

export function Register() {
	const [, setUserIdCookie] = useCookies(["user_id"]);
	const setUserIdStore = useUserStore((state) => state.setUserId);

	const [username, setUsername] = useState("");
	const [password, setPassword] = useState("");
	const [password2, setPassword2] = useState("");
	const [email, setEmail] = useState("");

	const router = useRouter();

	const { mutate: mutateRegister, isLoading: isLoadingRegister } =
		useRegister();
	const { mutate: mutateLogin, isLoading: isLoadingLogin } = useLogin();

	const registerForm = useForm(RegisterSchema, (values) => {
		mutateRegister(values, {
			onSuccess: () => {
				mutateLogin(
					{
						username: values.username,
						password: values.password,
					},
					{
						onSuccess: (mutationResp) => {
							if (mutationResp?.user_id && mutationResp?.user_id) {
								setUserIdCookie("user_id", mutationResp?.user_id);
								setUserIdStore(mutationResp?.user_id);
							}
							router.push("/");
						},
					}
				);
			},
		});
	});

	return (
		<form
			className="flex flex-col m-2 p-2"
			onSubmit={(e) => {
				e.preventDefault();
				registerForm.onSubmit({ username, email, password, password2 });
			}}
		>
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

			<Button isLoading={isLoadingRegister || isLoadingLogin} text="Register" />
		</form>
	);
}
