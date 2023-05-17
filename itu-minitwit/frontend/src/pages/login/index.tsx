import { Login } from "@/components/Auth/Login";
import DefaultLayout from "@/layouts/DefaultLayout";

export default function LoginPage() {
	return (
		<DefaultLayout>
			<div className="mt-10 bg-gray-300 w-96 rounded-md mx-auto">
				<Login />
			</div>
		</DefaultLayout>
	);
}
