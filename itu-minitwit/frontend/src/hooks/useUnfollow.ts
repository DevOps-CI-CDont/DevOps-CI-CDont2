import { makeFetch } from "@/lib/makeFetch";
import { queryClient } from "@/pages/_app";
import { FollowType } from "@/types/User.type";
import { useMutation } from "react-query";

export function useUnfollow() {
	return useMutation({
		mutationFn: async ({ userId, username }: FollowType) => {
			const headers = new Headers();
			headers.append("Authorization", userId);

			return await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/user/${username}/unfollow`,
				{
					method: "POST",
					headers,
					redirect: "follow",
				}
			);
		},
		onSuccess: () => {
			queryClient.invalidateQueries("isFollowing");
		},
	});
}
