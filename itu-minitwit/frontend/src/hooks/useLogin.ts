import { makeFetch } from "@/lib/makeFetch";
import { LoginSchemaType } from "@/types/Auth.type";
import { useMutation } from "react-query";

export function useLogin() {
	return useMutation({
		mutationFn: async ({ username, password }: LoginSchemaType) => {
			const formData = new FormData();
			formData.append("username", username);
			formData.append("password", password);

			return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/login`, {
				method: "POST",
				body: formData,
				redirect: "follow",
			}).then((res) => res.json());
		},
	});
}
