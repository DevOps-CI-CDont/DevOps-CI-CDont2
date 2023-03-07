export async function getPublicTweets() {
  const API_URL = process.env.NEXT_PUBLIC_API_URL as string
  console.log("getPublicTweets fetching from " + API_URL + "/public")
  const resp = await fetch(API_URL + "/public")
  console.log("getPublicTweets response: " + resp)
  const data = await resp.json()
  console.log("getPublicTweets data: " + data)
  return data
}
