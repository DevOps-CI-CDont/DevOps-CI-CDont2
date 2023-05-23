import { getAvatarUrl } from "@/lib/dicebear";
import useUserStore from "@/store/userStore";
import { Tweet as TweetType } from "@/types/tweet.type";
import clsx from "clsx";
import Image from "next/image";
import Link from "next/link";

export function Tweet({ text, author_name, CreatedAt, author_id }: TweetType) {
	const userId = useUserStore((state) => state.userId);

	const date = new Date(CreatedAt);

	const year = date.getFullYear();
	const month = date.getMonth() + 1;
	const day = date.getDate();
	const hours = date.getHours();
	const minutes = date.getMinutes();
	const seconds = date.getSeconds();

	const readableTimestamp = `${year}-${month < 10 ? "0" + month : month}-${
		day < 10 ? "0" + day : day
	} ${hours < 10 ? "0" + hours : hours}:${
		minutes < 10 ? "0" + minutes : minutes
	}:${seconds < 10 ? "0" + seconds : seconds}`;
	CreatedAt = readableTimestamp;

	return (
		<div
			className={clsx(
				"mx-4 my-4 py-2 px-2 bg-blue-500 flex rounded-md items-center dark:bg-blue-700",
				userId == author_id && "flex-row-reverse"
			)}
		>
			<div className="w-[4rem]">
				<Link href={`/user/${author_name}`}>
					<Image
						src={getAvatarUrl(author_name)}
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
						href={`/user/${author_name}`}
						prefetch={false}
						className="font-bold"
					>
						{author_name}
					</Link>{" "}
					at {CreatedAt}
				</span>
			</div>
		</div>
	);
}
