import { makeFetch } from "@/lib/makeFetch";
import { Tweet } from "@/types/tweet.type";
import { useQuery } from "react-query";

interface GetTweetsByUsernameProps {
	tweets: Tweet[] | null;
}

export function useGetTweetsByUsername(username: string) {
	return useQuery(["tweetsById", username], async () => {
		return await makeFetch<GetTweetsByUsernameProps>({
			url: `user/${username}`,
			method: "GET",
		});
	});
}
