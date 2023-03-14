import { ChangeEventHandler, InputHTMLAttributes } from "react";

interface InputProps {
	placeholder?: string;
	type: InputHTMLAttributes<HTMLInputElement>["type"];
	value: string;
	onChange: ChangeEventHandler<HTMLInputElement>;
	label: string;
}

export function Input({
	placeholder,
	type,
	value,
	onChange,
	label,
}: InputProps) {
	return (
		<div className="flex flex-col">
			<label className="text-sm font-bold">{label}</label>
			<input
				className="px-2 py-1 mb-4 mt-1 shadow-md rounded-md"
				placeholder={placeholder}
				type={type}
				value={value}
				onChange={onChange}
			/>
		</div>
	);
}
