import { useFollow } from "@/hooks/useFollow";
import { useGetIsFollowing } from "@/hooks/useGetIsFollowing";
import { useUnfollow } from "@/hooks/useUnfollow";
import { useRouter } from "next/router";
import { useCookies } from "react-cookie";

export function Follow() {
	const router = useRouter();

	const [userIdCookie] = useCookies(["user_id"]);

	const { username } = router.query;

	const {
		data: isFollowing,
		isLoading: isLoadingFollowing,
		isError: isFollowingError,
	} = useGetIsFollowing({
		username: username as string,
		userId: userIdCookie.user_id,
	});

	const { mutate: follow, isLoading: isLoadingFollow } = useFollow();
	const { mutate: unfollow, isLoading: isLoadingUnfollow } = useUnfollow();

	return (
		<button
			disabled={
				isLoadingFollowing ||
				isFollowingError ||
				isLoadingFollow ||
				isLoadingUnfollow
			}
			onClick={handleFollow}
			className="ml-2 font-bold px-2 py-1 border bg-blue-500 shadow-md text-white rounded-md disabled:bg-gray-200"
		>
			{isFollowing ? <>Unfollow</> : <>Follow</>}
		</button>
	);

	function handleFollow() {
		if (!isFollowing) {
			follow({
				username: username as string,
				userId: userIdCookie.user_id,
			});
		} else {
			unfollow({
				username: username as string,
				userId: userIdCookie.user_id,
			});
		}
	}
}
