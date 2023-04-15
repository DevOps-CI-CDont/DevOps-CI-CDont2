import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { CreateMessage } from "@/components/Message/CreateTweet";
import { useCookies } from "react-cookie";
import { useGetTimeline } from "@/hooks/useGetTimeline";

export default function MyTimelinePage() {
	const [userIdCookie] = useCookies(["user_id"]);
	const { data, isLoading } = useGetTimeline(userIdCookie as string);

	return (
		<DefaultLayout>
			<div className="wrapper mt-4">
				<h1 className="font-bold">My timeline</h1>
				<CreateMessage />
				{isLoading ? (
					<p>Loading...</p>
				) : (
					data?.tweets && <TweetContainer tweets={data.tweets} />
				)}
			</div>
		</DefaultLayout>
	);
}
