import { getAvatarUrl } from "@/lib/dicebear";
import useUserStore from "@/store/userStore";
import { Tweet as TweetType } from "@/types/tweet.type";
import clsx from "clsx";
import Image from "next/image";
import Link from "next/link";

export function Tweet({ text, author, pub_date, author_id }: TweetType) {
	const userId = useUserStore((state) => state.userId);

	return (
		<div
			className={clsx(
				"mx-4 my-4 py-2 px-2 bg-blue-500 flex rounded-md items-center",
				userId == author_id && "flex-row-reverse"
			)}
		>
			<div className="w-[4rem]">
				<Link href={`/user/${author.username}`}>
					<Image
						src={getAvatarUrl(author.username)}
						alt="avatar"
						width={60}
						height={60}
						className="rounded-full shadow-md"
					/>
				</Link>
			</div>
			<div
				className={clsx(
					"mx-4 my-2 text-white flex flex-col w-fit",
					userId == author_id && "text-right"
				)}
			>
				<span className="text-md">{text}</span>
				<span className="text-sm text-gray-200">
					Tweeted by{" "}
					<Link
						href={`/user/${author.username}`}
						prefetch={false}
						className="font-bold"
					>
						{author.username}
					</Link>{" "}
					at {new Date(pub_date * 1000).toLocaleString()}
				</span>
			</div>
		</div>
	);
}
