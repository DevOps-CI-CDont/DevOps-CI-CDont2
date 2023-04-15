import { Tweet } from "@/types/tweet.type";
import { useQuery } from "react-query";

interface GetPublicTimelineProps {
	tweets: Tweet[];
}

export function useGetPublicTimeline() {
	return useQuery(["publicTimeline"], async () => {
		return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/public`, {
			method: "GET",
			redirect: "follow",
		}).then((response) => response.json() as Promise<GetPublicTimelineProps>);
	});
}
