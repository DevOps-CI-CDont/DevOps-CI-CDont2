import { useQuery } from "react-query";

export function useGetTweetsByUsername(username: string) {
	return useQuery(["tweetsById", username], async () => {
		return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/user/${username}`, {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
			},
		}).then((response) => response.json());
	});
}
