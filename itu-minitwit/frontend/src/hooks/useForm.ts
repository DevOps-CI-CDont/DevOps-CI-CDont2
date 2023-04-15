import { z } from "zod";

export function useForm<TValues>(
	schema: z.Schema<TValues>,
	onSubmit: (values: TValues) => void
) {
	return {
		onSubmit: (values: TValues) => {
			const newValues = schema.parse(values);
			onSubmit(newValues);
		},
	};
}
