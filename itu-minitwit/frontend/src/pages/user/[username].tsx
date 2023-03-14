import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getIsFollowing } from "@/server/getIsFollowing";
import { getUserTweets } from "@/server/getUserTweets";
import { postIsFollowing } from "@/server/postFollow";
import { Tweet } from "@/types/tweet.type";
import { useEffect, useState } from "react";
import { useCookies } from "react-cookie";

interface UserProfileProps {
	username?: string;
	tweets?: Tweet[];
	isFollowing: boolean;
}

export default function UserProfilePage({
	username,
	tweets,
	isFollowing,
}: UserProfileProps) {
	const [userIdCookie] = useCookies(["user_id"]);
	const [followDisabled, setFollowDisabled] = useState<boolean>(false);
	const [isFollowingState, setIsFollowingState] =
		useState<boolean>(isFollowing);

	useEffect(() => {
		setFollowDisabled(userIdCookie.user_id == undefined ? false : true);
	}, [userIdCookie]);

	return (
		<DefaultLayout>
			<div className="mt-4">
				<div className="flex items-center justify-between">
					<h1 className="text-lg">{username}&apos;s profile</h1>
					<button
						disabled={!followDisabled}
						onClick={handleFollow}
						className="ml-2 font-bold px-2 py-1 border bg-blue-500 shadow-md text-white rounded-md disabled:bg-gray-200"
					>
						{isFollowingState ? "Unfollow" : "Follow"}
					</button>
				</div>
				<TweetContainer tweets={tweets} />
			</div>
		</DefaultLayout>
	);

	async function handleFollow() {
		try {
			await postIsFollowing({
				userId: userIdCookie.user_id,
				username: username,
				isFollowing: isFollowingState,
			});

			setIsFollowingState(!isFollowingState);
		} catch (e) {
			console.log(e);
		}
	}
}

export async function getServerSideProps(context: any) {
	const { username } = context.query;

	const cookies = String(context.req.headers?.cookie).split(";");

	let isFollowing = false;

	if (cookies) {
		const userId = cookies.find((cookie: string) => cookie.includes("user_id"));

		if (userId) {
			const userIdValue = userId.split("=")[1];

			isFollowing = await getIsFollowing({
				userId: userIdValue,
				username,
			});
		}
	}

	const messages = await getUserTweets({ username });

	if (!messages) {
		return {
			props: {
				username,
				tweets: [],
				isFollowing,
			},
		};
	}

	return {
		props: {
			username,
			tweets: messages.tweets,
			isFollowing,
		},
	};
}
