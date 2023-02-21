export async function getPublicTweets() {
  console.log("getPublicTweets fetching from proxy + : ", process.env.NEXT_PUBLIC_API_URL + "/public")
  const response = await fetch(
    process.env.NEXT_PUBLIC_PROXY_URL + "/" +
    process.env.NEXT_PUBLIC_API_URL + "/public")
  const data = await response.json()
  return data
}
