import { Tweet } from "@/types/tweet.type";
import { useQuery } from "react-query";

type GetTimelineProps = {
	tweets: Tweet[];
};

export function useGetTimeline(userId: string) {
	return useQuery(["timeline", userId], async () => {
		var myHeaders = new Headers();
		myHeaders.append("Cookie", `user_id=${userId}`);

		return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/mytimeline`, {
			method: "GET",
			headers: myHeaders,
			redirect: "follow",
		}).then((response) => response.json() as Promise<GetTimelineProps>);
	});
}
