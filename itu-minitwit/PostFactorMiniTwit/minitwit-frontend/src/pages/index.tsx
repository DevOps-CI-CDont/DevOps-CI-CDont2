import { CreateMessage } from "@/components/Message/CreateTweet";
import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getTimeline } from "@/server/getTimeline";
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
        {
          tweets && <TweetContainer tweets={tweets} />
        }
      </div>
    </DefaultLayout>
  );
}

export async function getServerSideProps(context: any) {
  try {
    const cookie = context.req.headers.cookie

    if(!cookie) {
      throw new Error("Not signed in")
    }

    const userId = cookie[8]

    const messages = await getTimeline(parseInt(userId));

    if (!messages.tweets) {
      return {
        props: {
          tweets: [],
        },
      }
    }

    return {
      props: {
        tweets: messages.tweets,
      },
    };
  } catch (e) {
    console.error(e);
    return {
      props: {
        tweets: [],
      },
    };
  }
}
