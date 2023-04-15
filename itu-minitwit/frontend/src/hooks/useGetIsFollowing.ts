import { FollowType } from "@/types/User.type";
import { useQuery } from "react-query";

export function useGetIsFollowing({ userId, username }: FollowType) {
	return useQuery(["isFollowing", userId, username], async () => {
		if (!userId || !username) return null;

		var myHeaders = new Headers();
		myHeaders.append("Cookie", `user_id=${userId}`);
		return await fetch(
			`${process.env.NEXT_PUBLIC_API_URL}/AmIFollowing/${username}`,
			{
				method: "GET",
				headers: myHeaders,
				redirect: "follow",
			}
		).then((response) => response.json());
	});
}
