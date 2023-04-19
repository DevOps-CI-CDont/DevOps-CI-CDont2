import { Loading } from "@/components/Loading";
import { TweetContainer } from "@/components/Message/TweetContainer";
import { useGetPublicTimeline } from "@/hooks/useGetPublicTimeline";
import DefaultLayout from "@/layouts/DefaultLayout";

export default function PublicTimelinePage() {
	const { data, isLoading } = useGetPublicTimeline();

	return (
		<DefaultLayout>
			<div className="mt-4">
				<h1>Public timeline</h1>
				{isLoading ? (
					<Loading />
				) : (
					data?.tweets && <TweetContainer tweets={data.tweets} />
				)}
			</div>
		</DefaultLayout>
	);
}
