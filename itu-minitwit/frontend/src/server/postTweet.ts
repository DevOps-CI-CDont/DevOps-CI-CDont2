export interface PostTweetProps {
  message: string;
  userId: string;
}

export async function postTweet({ message, userId }: PostTweetProps) {
  let formData = new FormData();
  formData.append("text", message);
  return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/add_message`, {
    mode: "no-cors",
    method: "POST",
    cache: "no-cache",
    headers: {
      Cookie: `user_id=${userId}`,
      "Content-Type": "application/json",
      origin: "http://localhost:3000",
    },
    credentials: "include",
    redirect: "follow",
    body: formData,
  }).then((response) => response.json());
}
