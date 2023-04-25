import { RegisterSchemaType } from "@/types/Auth.type";
import { useMutation } from "react-query";

export function useRegister() {
	return useMutation({
		mutationFn: async ({
			username,
			email,
			password,
			password2,
		}: RegisterSchemaType) => {
			const formData = new FormData();
			formData.append("username", username);
			formData.append("email", email);
			formData.append("password", password);
			formData.append("password2", password2);

			return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/register`, {
				method: "POST",
				body: formData,
				redirect: "follow",
			}).then((res) => res.json());
		},
	});
}
