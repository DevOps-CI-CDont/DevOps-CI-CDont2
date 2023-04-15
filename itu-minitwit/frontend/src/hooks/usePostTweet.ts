import { useMutation } from "react-query";
import { queryClient } from "@/pages/_app";
import { PostTweetSchemaType } from "@/types/tweet.type";

export function usePostTweet() {
	const postTweetMutation = useMutation({
		mutationFn: async ({ message, userId }: PostTweetSchemaType) => {
			let formData = new FormData();
			formData.append("text", message);
			return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/add_message`, {
				mode: "no-cors",
				method: "POST",
				cache: "no-cache",
				headers: {
					Cookie: `user_id=${userId}`,
					"Content-Type": "application/json",
					origin: "http://localhost:3000",
				},
				credentials: "include",
				redirect: "follow",
				body: formData,
			});
		},
		onSuccess: () => {
			queryClient.invalidateQueries(["timeline", "publicTimeline"]);
		},
	});

	return postTweetMutation;
}
