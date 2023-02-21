export async function getPublicTweets() {
  const PROXY_URL = process.env.NEXT_PUBLIC_PROXY_URL as string
  const API_URL = process.env.NEXT_PUBLIC_API_URL as string

  console.log("getPublicTweets fetching from " + PROXY_URL + API_URL + "/public")

  const response = await fetch(
    PROXY_URL + "/" +
    API_URL + "/public")
  const data = await response.json()
  return data
}
