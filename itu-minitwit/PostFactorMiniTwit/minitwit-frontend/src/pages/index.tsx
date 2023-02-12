import { CreateMessage } from "@/components/Message/CreateTweet";
import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getPublicTweets } from "@/server/getPublicTweets";
import { Tweet } from "@/types/tweet.type";

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

export async function getServerSideProps() {
  try {
    const messages = await getPublicTweets();

    if (!messages.tweets) {
      return {
        props: {
          tweets: [],
        },
      };
    }

    return {
      props: {
        tweets: messages.tweets,
      },
    };
  } catch (e) {
    return {
      props: {
        tweets: [],
      },
    };
  }
}
