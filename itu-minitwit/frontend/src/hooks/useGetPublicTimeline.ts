import { makeFetch } from "@/lib/makeFetch";
import { Tweet } from "@/types/tweet.type";
import { useQuery } from "react-query";

interface GetPublicTimelineProps {
	tweets: Tweet[] | null;
}

export function useGetPublicTimeline() {
	return useQuery(["publicTimeline"], async () => {
		return await makeFetch<GetPublicTimelineProps>({
			url: "public",
			method: "GET",
		});
	});
}
