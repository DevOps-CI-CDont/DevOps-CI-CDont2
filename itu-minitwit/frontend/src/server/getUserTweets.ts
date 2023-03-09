interface GetUserTweetsProps {
  username: string;
}

export async function getUserTweets({ username }: GetUserTweetsProps) {
  return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/user/${username}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  }).then((response) => response.json());
}
