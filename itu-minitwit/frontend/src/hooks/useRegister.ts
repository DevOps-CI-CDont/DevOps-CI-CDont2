import { postLogin } from "@/server/postLogin";
import { RegisterSchemaType } from "@/types/Auth.type";
import { useRouter } from "next/router";
import { useMutation } from "react-query";
import { useLogin } from "./useLogin";

export function useRegister() {
	const router = useRouter();
	const mutateLogin = useLogin();

	return useMutation({
		mutationFn: async ({
			username,
			email,
			password,
			password2,
		}: RegisterSchemaType) => {
			var formdata = new FormData();
			formdata.append("username", username);
			formdata.append("email", email);
			formdata.append("password", password);
			formdata.append("password2", password2);

			return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/register`, {
				method: "POST",
				cache: "no-cache",
				redirect: "follow",
				body: formdata,
			}).then((response) => response.json());
		},
	});
}
