import { makeFetch } from "@/lib/makeFetch";
import { Tweet } from "@/types/tweet.type";
import { useQuery } from "react-query";

type GetTimelineProps = {
	tweets: Tweet[] | null;
};

export function useGetTimeline(userId: string) {
	return useQuery(["mytimeline", userId], async () => {
		return await makeFetch<GetTimelineProps>({
			method: "GET",
			url: "mytimeline",
			userId,
		});
	});
}
