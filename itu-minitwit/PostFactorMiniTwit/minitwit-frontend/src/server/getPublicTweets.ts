export async function getPublicTweets() {
  return await fetch(
    `${process.env.NEXT_PUBLIC_CORS_ORIGIN}/${process.env.NEXT_PUBLIC_API_URL}/public`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        origin: "http://localhost:3000",
      },
      redirect: "follow",
    }
  ).then((response) => response.json());
}
