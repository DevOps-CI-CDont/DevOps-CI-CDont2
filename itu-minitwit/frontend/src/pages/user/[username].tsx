import { Loading } from "@/components/Loading";
import { TweetContainer } from "@/components/Message/TweetContainer";
import { useFollow } from "@/hooks/useFollow";
import { useGetIsFollowing } from "@/hooks/useGetIsFollowing";
import { useGetTweetsByUsername } from "@/hooks/useGetTweetsByUserId";
import { useUnfollow } from "@/hooks/useUnfollow";
import DefaultLayout from "@/layouts/DefaultLayout";
import { FollowSchema } from "@/types/User.type";

import { useRouter } from "next/router";
import { useCookies } from "react-cookie";

export default function UserProfilePage() {
	const router = useRouter();

	const [userIdCookie] = useCookies(["user_id"]);

	const { username } = router.query;

	const { data: userTweets, isLoading } = useGetTweetsByUsername(
		username as string
	);

	const { data: isFollowing, isLoading: isLoadingFollowing } =
		useGetIsFollowing({
			username: username as string,
			userId: userIdCookie.user_id as string,
		});

	const { mutate: mutateFollow, isLoading: isLoadingFollow } = useFollow();
	const { mutate: mutateUnfollow, isLoading: isLoadingUnfollow } =
		useUnfollow();

	return (
		<DefaultLayout>
			<div className="mt-4">
				<div className="flex items-center justify-between">
					<h1 className="text-lg">{username}&apos;s profile</h1>
					<button
						disabled={
							isLoadingFollowing || isLoadingFollow || isLoadingUnfollow
						}
						onClick={handleFollow}
						className="ml-2 font-bold px-2 py-1 border bg-blue-500 shadow-md text-white rounded-md disabled:bg-gray-200"
					>
						{isFollowing ? "Unfollow" : "Follow"}
					</button>
				</div>
				{isLoading ? (
					<Loading />
				) : (
					<TweetContainer tweets={userTweets?.tweets} />
				)}
			</div>
		</DefaultLayout>
	);

	async function handleFollow() {
		const userfollowRequest = {
			username: username as string,
			userId: userIdCookie?.user_id as string,
		};

		if (!FollowSchema.parse(userfollowRequest)) {
			return alert("Something went wrong");
		}
		if (isFollowing) {
			mutateUnfollow(userfollowRequest, {
				onSuccess: (data) => {
					// console.log(data);
				},
				onError: (error) => {
					console.error(error);
				},
			});
		} else {
			mutateFollow(userfollowRequest, {
				onSuccess: (data) => {
					// console.log(data);
				},
				onError: (error) => {
					console.error(error);
				},
			});
		}
	}
}
