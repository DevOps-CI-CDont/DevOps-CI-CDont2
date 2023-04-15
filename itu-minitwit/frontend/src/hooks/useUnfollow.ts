import { queryClient } from "@/pages/_app";
import { FollowType } from "@/types/User.type";
import { useMutation } from "react-query";

export function useUnfollow() {
	return useMutation({
		mutationFn: async ({ userId, username }: FollowType) => {
			var myHeaders = new Headers();
			myHeaders.append("Cookie", `user_id=${userId}`);

			return await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/user/${username}/unfollow`,
				{
					method: "POST",
					headers: myHeaders,
					redirect: "follow",
				}
			).then((response) => response.json());
		},
		onSuccess: () => {
			queryClient.invalidateQueries("isFollowing");
		},
	});
}
