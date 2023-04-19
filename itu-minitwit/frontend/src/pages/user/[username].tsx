import { Follow } from "@/components/Follow";
import { Loading } from "@/components/Loading";
import { TweetContainer } from "@/components/Message/TweetContainer";
import { useGetTweetsByUsername } from "@/hooks/useGetTweetsByUserId";
import DefaultLayout from "@/layouts/DefaultLayout";
import { useRouter } from "next/router";

export default function UserProfilePage() {
	const router = useRouter();

	const { username } = router.query;

	const { data: userTweets, isLoading } = useGetTweetsByUsername(
		username as string
	);

	return (
		<DefaultLayout>
			<div className="mt-4">
				<div className="flex items-center justify-between">
					<h1 className="text-lg">{username}&apos;s profile</h1>
					<Follow />
				</div>
				{isLoading ? (
					<Loading />
				) : (
					<TweetContainer tweets={userTweets?.tweets} />
				)}
			</div>
		</DefaultLayout>
	);
}
