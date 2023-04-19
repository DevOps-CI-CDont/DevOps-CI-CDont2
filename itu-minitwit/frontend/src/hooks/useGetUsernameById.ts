import { makeFetch } from "@/lib/makeFetch";
import { useQuery } from "react-query";

type GetUsernameByIdError = {
	error: string;
};

export function useGetUsernameById(userId: string | number) {
	return useQuery(["username", userId], async () => {
		if (!userId) return null;

		return await makeFetch<GetUsernameByIdError | string>({
			url: `getUserNameById?id=${userId}`,
			method: "GET",
		});
	});
}
