import axios from "axios";

export async function getUsernameById(userId: string | number) {
	return await axios
		.get(
			`${
				process.env.NEXT_PUBLIC_API_URL
			}/getUserNameById?id=${userId.toString()}`
		)
		.then((response) => response.data);
}
