export async function getPublicTweets() {
  const API_URL = process.env.NEXT_PUBLIC_API_URL as string;

  const resp = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/public`);

  const data = await resp.json();

  return data;
}
