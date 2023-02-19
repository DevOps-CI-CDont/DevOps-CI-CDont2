interface GetIsFollowingProps {
    userId: string;
    username: string;
}

export async function getIsFollowing({userId, username}: GetIsFollowingProps) {
    return await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/AmIFollowing/${username}`,
        {
          method: "GET",
          mode: "no-cors",
          cache: "no-cache",
          headers: {
            "Cookie": `user_id=${userId}`,
            "Content-Type": "application/json",
            origin: "http://localhost:3000",
          },
          redirect: "follow",
        }
      ).then((response) => response.json());
}