import { z } from "zod";

const LoginSchema = z.object({
	username: z.string(),
	password: z.string().min(8),
});

type LoginSchemaType = z.infer<typeof LoginSchema>;

const RegisterSchema = z.object({
	username: z.string().min(3),
	email: z.string().email(),
	password: z.string().min(8),
	password2: z.string().min(8),
});
type RegisterSchemaType = z.infer<typeof RegisterSchema>;

export { LoginSchema, RegisterSchema };
export type { LoginSchemaType, RegisterSchemaType };
