import { useQuery } from "react-query";

export function useGetUsernameById(userId: string | number) {
	return useQuery(["username", userId], async () => {
		return await fetch(
			`${
				process.env.NEXT_PUBLIC_API_URL
			}/getUserNameById?id=${userId.toString()}`
		).then((response) => response.json());
	});
}
