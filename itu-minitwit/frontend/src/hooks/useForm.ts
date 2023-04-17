import { z } from "zod";

export function useForm<TValues>(
	schema: z.Schema<TValues>,
	onSubmit: (values: TValues) => void
) {
	return {
		onSubmit: (values: TValues) => {
			const parsedValues = schema.safeParse(values);

			if (parsedValues.success) {
				onSubmit(parsedValues.data);
			} else {
				alert("Something went wrong \n" + parsedValues.error.message);
			}
		},
	};
}
