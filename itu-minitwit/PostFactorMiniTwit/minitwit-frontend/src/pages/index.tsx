import { CreateMessage } from "@/components/Message/CreateTweet";
import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getPublicTweets } from "@/server/getPublicTweets";
import { Tweet } from "@/types/tweet.type";
import * as cookie from "cookie";

interface MyTimelinePageProps {
  tweets?: Tweet[];
}

export default function MyTimelinePage({ tweets }: MyTimelinePageProps) {
  return (
    <DefaultLayout>
      <div className='wrapper mt-4'>
        <h1 className='font-bold'>My timeline</h1>
        <CreateMessage />
        <TweetContainer tweets={tweets} />
      </div>
    </DefaultLayout>
  );
}

export async function getServerSideProps(context: any) {
  try {
    const c = cookie.parse(context.req.headers.cookie);

    const parsed = JSON.parse(c.session);

    const messages = await getPublicTweets();

    if (!messages.tweets) {
      return {
        props: {
          tweets: [],
        },
      };
    }

    const filteredMessages = messages.tweets.filter(
      (tweet: Tweet) => tweet.author_id === parseInt(parsed.user)
    );

    return {
      props: {
        tweets: filteredMessages,
      },
    };
  } catch (e) {
    console.log(e);
    return {
      props: {
        tweets: [],
      },
    };
  }
}
