export async function getPublicTweets() {
  const PROXY_URL = process.env.NEXT_PUBLIC_PROXY_URL as string
  const API_URL = process.env.NEXT_PUBLIC_API_URL as string
  console.log("getPublicTweets fetching from " + API_URL + "/public")
  const resp = await fetch(process.env.NEXT_PUBLIC_API_URL + "/public")
  console.log("getPublicTweets response: " + resp)
  const data = await resp.json()
  console.log("getPublicTweets data: " + data)
  return data
}
