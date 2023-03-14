import { getIsFollowing } from "./getIsFollowing";

interface PostWillFollow extends GetIsFollowingProps {
  isFollowing: boolean;
}

export async function postIsFollowing({
  userId,
  username,
  isFollowing,
}: PostWillFollow) {
  if (!userId || !username) {
    return;
  }
  if (isFollowing) {
    return await unFollow({ userId, username });
  } else {
    return await follow({ userId, username });
  }
}

interface GetIsFollowingProps {
  userId?: string;
  username?: string;
}

async function unFollow({ userId, username }: GetIsFollowingProps) {
  if (!username || !userId) {
    return;
  }

  try {
    var myHeaders = new Headers();
    myHeaders.append("Cookie", "user_id=5");

    var formdata = new FormData();

    fetch(`${process.env.NEXT_PUBLIC_API_URL}/user/${username}/unfollow`, {
      method: "POST",
      headers: myHeaders,
      body: formdata,
      redirect: "follow",
    }).then((response) => response.text());
  } catch (e) {
    console.error(e);
  }
}

async function follow({ userId, username }: GetIsFollowingProps) {
  if (!username || !userId) {
    return;
  }
  try {
    var myHeaders = new Headers();
    myHeaders.append("Cookie", "user_id=5");

    var formdata = new FormData();

    fetch(`${process.env.NEXT_PUBLIC_API_URL}/user/${username}/follow`, {
      method: "POST",
      headers: myHeaders,
      body: formdata,
      redirect: "follow",
    }).then((response) => response.text());
  } catch (e) {
    console.error(e);
  }
}
