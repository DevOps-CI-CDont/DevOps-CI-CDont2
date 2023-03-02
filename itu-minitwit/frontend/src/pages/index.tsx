import { CreateMessage } from "@/components/Message/CreateTweet";
import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getPublicTweets } from "@/server/getPublicTweets";
import { Tweet } from "@/types/tweet.type";
import { useCookies } from "react-cookie";
import { useEffect, useState } from "react";

export default function MyTimelinePage() {
  const [cookies] = useCookies(["user_id"]);

  const [isLogged, setIsLogged] = useState(false);
  const [tweets, setTweets] = useState<Tweet[]>([]);

  useEffect(() => {
    getPublicTweets().then((res) => setTweets(res.tweets));
  }, []);

  useEffect(() => {
    if (cookies.user_id) {
      setIsLogged(true);
    }
  }, [cookies]);

  if (!isLogged) return <></>;

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

async function getServerSideProps(context: any) {}
//   try {
//     const cookie = context.req.headers.cookie

//     if(!cookie) {
//       throw new Error("Not signed in")
//     }

//     const userId = cookie[8]

//     const messages = await getTimeline(parseInt(userId));

//     if (!messages.tweets) {
//       return {
//         props: {
//           tweets: [],
//         },
//       }
//     }

//     return {
//       props: {
//         tweets: messages.tweets,
//       },
//     };
//   } catch (e) {
//     console.error(e);
//     return {
//       props: {
//         tweets: [],
//       },
//     };
//   }
// }
