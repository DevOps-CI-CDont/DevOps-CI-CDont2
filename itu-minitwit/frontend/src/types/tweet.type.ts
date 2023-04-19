import { z } from "zod";

export type Tweet = {
	CreatedAt: string;
	DeledtedAt: string | null;
	UpdatedAt: string | null;
	ID: number;
	author_id: number;
	author_name: string;
	text: string;
	pub_date: number;
	flagged: number;
};

const postTweetSchema = z.object({
	message: z.string(),
	userId: z.string(),
});

type PostTweetSchemaType = z.infer<typeof postTweetSchema>;

export { postTweetSchema };
export type { PostTweetSchemaType };
