export async function getUsernameById(userId: string | number) {
	return await fetch(
		`${process.env.NEXT_PUBLIC_API_URL}/getUserNameById?id=${userId.toString()}`
	).then((response) => response.json());
}
