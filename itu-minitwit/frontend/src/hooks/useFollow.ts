import { queryClient } from "@/pages/_app";
import { FollowType } from "@/types/User.type";
import { useMutation } from "react-query";
import axios from "axios";

export function useFollow() {
	return useMutation({
		mutationFn: async ({ username, userId }: FollowType) => {
			var myHeaders = new Headers();
			myHeaders.append("Cookie", `user_id=${userId}`);

			return await fetch(
				`${process.env.NEXT_PUBLIC_API_URL}/user/${username}/follow`,
				{
					credentials: "same-origin",
					method: "POST",
					headers: {
						"Content-Type": "application/json",
						Cookie: `user_id=${userId}`,
					},
					redirect: "follow",
				}
			).then((response) => response.json());
		},
		onSuccess: () => {
			queryClient.invalidateQueries("isFollowing");
		},
	});
}
