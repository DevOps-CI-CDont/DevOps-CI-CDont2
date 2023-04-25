import { z } from "zod";

const FollowSchema = z.object({
	username: z.string(),
	userId: z.string(),
});

type FollowType = z.infer<typeof FollowSchema>;

export { FollowSchema };
export type { FollowType };
