import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getPublicTweets } from "@/server/getPublicTweets";
import { Tweet } from "@/types/tweet.type";
import { useEffect, useState } from "react";

interface HomePageProps {
  tweets: Tweet[];
}

export default function PublicTimelinePage({ tweets }: HomePageProps) {
  // useEffect(() => {
  //   getPublicTweets().then((res) => setTweets(res.tweets));
  // }, []);

  return (
    <DefaultLayout>
      <div className='mt-4'>
        <h1>Public timeline</h1>
        <TweetContainer tweets={tweets} />
      </div>
    </DefaultLayout>
  );
}

export async function getServerSideProps() {
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
}
