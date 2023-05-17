import { usePostTweet } from "@/hooks/usePostTweet";
import { postTweetSchema } from "@/types/tweet.type";
import { useEffect, useState } from "react";
import { useCookies } from "react-cookie";
import { Button } from "../Input/Button";

export function CreateMessage() {
	const [message, setMessage] = useState("");

	const [userIdCookie] = useCookies(["user_id"]);
	const [hasCookie, setHasCookie] = useState(false);

	const postTweetMutation = usePostTweet();

	useEffect(() => {
		if (userIdCookie.user_id) {
			setHasCookie(true);
		} else {
			setHasCookie(false);
		}
	}, [userIdCookie]);

	if (!hasCookie) return <></>;

	return (
		<div className="w-full mt-2 bg-gray-300 dark:bg-slate-900 dark:text-slate-100 shadow-md px-1 py-1 rounded-md">
			<div className="mx-4 my-2">
				<span>What&apos;s on your mind?</span>
				<form className="flex flex-row mt-2" onSubmit={(e) => handleSubmit(e)}>
					<input
						type="text"
						value={message}
						disabled={!hasCookie || postTweetMutation.isLoading}
						onChange={(e) => setMessage(e.target.value)}
						className="px-2 py-1 mr-2 rounded-md w-full border shadown-md dark:border-slate-900"
						placeholder="write here..."
					/>
					<Button isLoading={postTweetMutation.isLoading} text={"Share"} />
				</form>
			</div>
		</div>
	);

	async function handleSubmit(e: any) {
		e.preventDefault();

		if (!hasCookie) {
			return alert("You need to be logged in to post a tweet");
		}

		const tweet = {
			message,
			userId: userIdCookie.user_id,
		};

		if (!postTweetSchema.parse(tweet)) {
			return alert("Invalid tweet");
		} else {
			postTweetMutation.mutate(tweet, {
				onSuccess: () => {
					setMessage("");
				},
			});
		}
	}
}
