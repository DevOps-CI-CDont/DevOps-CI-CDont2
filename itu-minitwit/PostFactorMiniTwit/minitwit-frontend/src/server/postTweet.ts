export async function postTweet() {
  return await fetch(
    `${process.env.NEXT_PUBLIC_CORS_ORIGIN}/${process.env.NEXT_PUBLIC_API_URL}/add_message`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        origin: "http://localhost:3000",
      },
      credentials: "same-origin",
      redirect: "follow",
      body: JSON.stringify({
        text: "hello",
      }),
    }
  ).then((response) => response.json());
}
