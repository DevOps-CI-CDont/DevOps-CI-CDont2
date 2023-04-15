import { LoginSchemaType } from "@/types/Auth.type";
import { useMutation } from "react-query";

export function useLogin() {
	return useMutation({
		mutationFn: async ({ username, password }: LoginSchemaType) => {
			var formdata = new FormData();
			formdata.append("username", username);
			formdata.append("password", password);

			return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/login`, {
				method: "POST",
				body: formdata,
			}).then((response) => response?.json());
		},
	});
}
