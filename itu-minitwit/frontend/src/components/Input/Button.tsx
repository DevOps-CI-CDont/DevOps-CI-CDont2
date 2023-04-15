import { Loading } from "../Loading";

interface ButtonProps {
	text: string;
	isLoading?: boolean;
}

export function Button({ text, isLoading = false }: ButtonProps) {
	return (
		<button
			disabled={isLoading}
			className="border bg-white rounded-md shadow-md my-4"
			type="submit"
		>
			{isLoading ? <Loading /> : <>{text}</>}
		</button>
	);
}
