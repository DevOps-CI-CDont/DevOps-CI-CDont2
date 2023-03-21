export async function getPublicTweets() {
	const resp = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/public`);

	const data = await resp.json();

	return data;
}
