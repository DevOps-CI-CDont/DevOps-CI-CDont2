import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { Tweet } from "@/types/tweet.type";
import { useCookies } from "react-cookie";
import { getTimeline } from "@/server/getTimeline";
import { CreateMessage } from "@/components/Message/CreateTweet";
import { useGetTimeline } from "@/hooks/useGetTimeline";

interface MyTimeLinePageProps {
  tweets: Tweet[];
}

export default function MyTimelinePage({ tweets }: MyTimeLinePageProps) {
  const [cookies] = useCookies(["user_id"]);

  return (
    <DefaultLayout>
      <div className='wrapper mt-4'>
        <h1 className='font-bold'>My timeline</h1>
        <CreateMessage />
        {tweets && <TweetContainer tweets={tweets} />}
      </div>
    </DefaultLayout>
  );
}

export async function getServerSideProps(context: any) {
  try {
    const cookie = context.req.headers.cookie;

    const cookies = String(context.req.headers?.cookie).split(";");

    let messages = [];

    if (cookies) {
      const userId = cookies.find((cookie: string) =>
        cookie.includes("user_id")
      );

      if (userId) {
        const userIdValue = userId.split("=")[1];

        messages = await getTimeline({ userId: userIdValue });
      }
    }

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
    console.error(e);
    return {
      props: {
        tweets: [],
      },
    };
  }
}
