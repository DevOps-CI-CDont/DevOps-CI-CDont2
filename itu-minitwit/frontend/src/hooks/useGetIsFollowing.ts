import { makeFetch } from "@/lib/makeFetch";
import { FollowType } from "@/types/User.type";
import { useQuery } from "react-query";

export function useGetIsFollowing({ userId, username }: FollowType) {
	return useQuery(["isFollowing", userId, username], async () => {
		if (!userId || !username) return null;

		return await makeFetch({
			url: `AmIFollowing/${username}`,
			method: "GET",
			userId,
		});
	});
}
