export async function postTweet(message: string) {
  let formData = new FormData()
  formData.append("text", message)
  console.log("formdata text", formData.get("text"))
  return await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/add_message`,
    {
      mode: "no-cors",
      method: "POST",
      cache: "no-cache",
      headers: {
        "Content-Type": "application/json",
        origin: "http://localhost:3000",
      },
      credentials: "include",
      redirect: "follow",
      body: formData
    }
  ).then((response) => response.json());
}
