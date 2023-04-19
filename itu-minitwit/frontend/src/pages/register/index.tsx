import { Register } from "@/components/Auth/Register";
import DefaultLayout from "@/layouts/DefaultLayout";

export default function RegisterPage() {
	return (
		<DefaultLayout>
			<div className="mt-10 bg-gray-300 w-96 rounded-md mx-auto">
				<Register />
			</div>
		</DefaultLayout>
	);
}
