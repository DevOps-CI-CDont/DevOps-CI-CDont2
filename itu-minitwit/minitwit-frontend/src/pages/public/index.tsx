import { TweetContainer } from "@/components/Message/TweetContainer";
import DefaultLayout from "@/layouts/DefaultLayout";
import { getPublicTweets } from "@/server/getPublicTweets";
import { Tweet } from "@/types/tweet.type";

interface HomePageProps {
  tweets: Tweet[];
}

export default function PublicTimelinePage() {
  const [tweets, setTweets] = useState<Tweet[]>([]);

  useEffect(() => {
    getPublicTweets().then((res) => setTweets(res.tweets));

  }, []);
  
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
    console.error(e);

    return {
      props: {
        tweets: [],
      },
    };
  }
}
function useEffect(arg0: () => void, arg1: never[]) {
  throw new Error("Function not implemented.");
}

function useState<T>(arg0: never[]): [any, any] {
  throw new Error("Function not implemented.");
}

