import { getIsFollowing } from "./getIsFollowing";

interface PostWillFollow extends GetIsFollowingProps {
  isFollowing: boolean
}

export async function postIsFollowing({userId, username, isFollowing}: PostWillFollow) {
  if(!userId || !username) {
    return;
  }


    if(isFollowing) {
        return await unFollow({userId, username})
    } else {
        return await follow({userId, username})
    }
}

interface GetIsFollowingProps {
  userId?: string;
  username?: string;
}

async function unFollow({userId, username}: GetIsFollowingProps) {
  if(!username || !userId) {
    return
  }

  try {
    return await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/user/${username}/unfollow`,
      {
        method: "POST",
        mode: 'no-cors',
        headers: {
          "Cookie": `user_id=${userId}`,
          "Content-Type": "application/json",
          origin: "http://localhost:3000",
        },
        credentials: "include",
        redirect: "follow",
      }
    ).then((response) => response.text());
  } catch(e) {
    console.error(e)
  }

    
}

async function follow({userId, username}: GetIsFollowingProps) {
  if(!username || !userId) {
    return
  }
  try {
    return await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/user/${username}/follow`,
      {
        method: "POST",
        mode: 'no-cors',
        headers: {
          "Cookie": `user_id=${userId}`,
          "Content-Type": "application/json",
          origin: "http://localhost:3000",
        },
        credentials: "include",
        redirect: "follow",
      }
    ).then((response) => response.text());
  } catch(e) {
    console.error(e)
  }

    
}