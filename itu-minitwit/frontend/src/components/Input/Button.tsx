import { Loading } from "../Loading";

interface ButtonProps {
	text: string;
	isLoading?: boolean;
}

export function Button({ text, isLoading = false }: ButtonProps) {
	return (
		<button
			disabled={isLoading}
			className="border bg-white dark:bg-blue-700 dark:border-slate-900 rounded-md shadow-md py-2 px-4"
			type="submit"
		>
			{isLoading ? <Loading /> : <>{text}</>}
		</button>
	);
}
