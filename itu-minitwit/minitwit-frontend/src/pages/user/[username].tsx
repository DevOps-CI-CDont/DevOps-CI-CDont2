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
  isFollowing
}: UserProfileProps) {
  const [userIdCookie] = useCookies(['user_id'])
  const [follow, setFollow] = useState<boolean>(isFollowing)


  return (
    <DefaultLayout>
      <div className='mt-4'>
        <div className="flex items-center justify-between">
          <h1 className="text-lg">{username}&apos;s profile</h1>
          <button onClick={handleFollow} className="ml-2 font-bold px-2 py-1 border bg-blue-500 shadow-md text-white rounded-md">{
            follow ? 'Unfollow' : 'Follow'
          }</button>
        </div>
        <TweetContainer tweets={tweets} />
      </div>
    </DefaultLayout>
  );

  async function handleFollow() {
    try {
      await postIsFollowing({userId: userIdCookie.user_id, username: username, isFollowing: follow})

    setFollow(!follow)
    } catch(e) {}

  }
}



export async function getServerSideProps(context: any) {
  try {
    const { username } = context.query;

    const cookie = context.req.headers.cookie

    if(!cookie) {
      throw new Error("Not signed in")
    }

    const userId = cookie[8]

    const messages = await getPublicTweets();

    const isFollowing = await getIsFollowing({userId, username});

    if (!messages.tweets) {
      throw new Error("No tweets found");
    }

    const filteredTweets = messages.tweets.filter(
      (tweet: Tweet) => tweet.author.username === username
    );

    return {
      props: {
        username,
        tweets: filteredTweets,
        isFollowing: isFollowing
      },
    };
  } catch (e) {
    return {
      notFound: true,
    };
  }
}
