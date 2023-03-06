import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getIsFollowing } from "@/server/getIsFollowing";
import { getPublicTweets } from "@/server/getPublicTweets";
import { postIsFollowing } from "@/server/postFollow";
import { Tweet } from "@/types/tweet.type";
import { useState } from "react";
import { useCookies } from "react-cookie";

interface UserProfileProps {
  username?: string;
  tweets?: Tweet[];
  isFollowing: boolean;
}

export default function UserProfilePage({
  username,
  tweets,
  isFollowing,
}: UserProfileProps) {
  const [userIdCookie] = useCookies(["user_id"]);
  const [canFollow] = useState<boolean>(userIdCookie.user_id ? true : false);
  const [follow, setFollow] = useState<boolean>(isFollowing);

  return (
    <DefaultLayout>
      <div className='mt-4'>
        <div className='flex items-center justify-between'>
          <h1 className='text-lg'>{username}&apos;s profile</h1>
          <button
            disabled={!canFollow}
            onClick={handleFollow}
            className='ml-2 font-bold px-2 py-1 border bg-blue-500 shadow-md text-white rounded-md disabled:bg-gray-200'>
            {follow ? "Unfollow" : "Follow"}
          </button>
        </div>
        <TweetContainer tweets={tweets} />
      </div>
    </DefaultLayout>
  );

  async function handleFollow() {
    try {
      await postIsFollowing({
        userId: userIdCookie.user_id,
        username: username,
        isFollowing: follow,
      });

      setFollow(!follow);
    } catch (e) {}
  }
}

export async function getServerSideProps(context: any) {
  const { username } = context.query;

  const cookie = context.req.headers.cookie;
  const messages = await getPublicTweets();

  if (!messages) {
    return {
      props: {
        username,
        tweets: [],
        isFollowing: false,
      },
    };
  }

  const filteredTweets = messages.tweets.filter(
    (tweet: Tweet) => tweet.author.username === username
  );

  const userId = cookie && cookie[8];

  if (!userId) {
    return {
      props: {
        username,
        tweets: messages.tweets,
        isFollowing: false,
      },
    };
  }

  const isFollowing = await getIsFollowing({ userId, username });

  return {
    props: {
      username,
      tweets: filteredTweets,
      isFollowing: isFollowing,
    },
  };
}
