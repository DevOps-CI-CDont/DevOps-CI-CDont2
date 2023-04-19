import { useMutation } from "react-query";
import { queryClient } from "@/pages/_app";
import { PostTweetSchemaType } from "@/types/tweet.type";

export function usePostTweet() {
	const postTweetMutation = useMutation({
		mutationFn: async ({ message, userId }: PostTweetSchemaType) => {
			const formData = new FormData();
			formData.append("text", message);

			const headers = new Headers();
			headers.append("Authorization", userId);

			return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/add_message`, {
				method: "POST",
				body: formData,
				headers,
				redirect: "follow",
			}).then((res) => res.json());
		},
		onSuccess: () => {
			queryClient.invalidateQueries("mytimeline");
			queryClient.invalidateQueries("publicTimeline");
		},
	});

	return postTweetMutation;
}
