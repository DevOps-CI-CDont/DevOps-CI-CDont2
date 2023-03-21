import { Tweet as TweetType } from "@/types/tweet.type";
import { Tweet } from "./Tweet";

interface TweetContainerProps {
	tweets?: TweetType[];
}

export function TweetContainer({ tweets }: TweetContainerProps) {
	return (
		<div className="w-full mt-2 bg-gray-200 shadow-md px-1 py-1 mb-4 rounded-md">
			{!tweets ? (
				<span>No tweets found</span>
			) : (
				tweets?.map((tweets) => {
					return <Tweet {...tweets} key={tweets.ID} />;
				})
			)}
		</div>
	);
}
