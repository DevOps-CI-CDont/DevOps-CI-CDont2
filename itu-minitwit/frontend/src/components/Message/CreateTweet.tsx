import { getPublicTweets } from "@/server/getPublicTweets";
import { postTweet } from "@/server/postTweet";
import { Tweet } from "@/types/tweet.type";
import { useEffect, useState } from "react";
import { useCookies } from "react-cookie";

interface CreateMessageProps {
  setTweets: (tweets: Tweet[]) => void;
}


export function CreateMessage({ setTweets }: CreateMessageProps) {
  const [message, setMessage] = useState("");
  const [cookies] = useCookies(["user_id"]);

  const [hasCookie, setHasCookie] = useState(false);

  useEffect(() => {
    if (cookies.user_id) {
      setHasCookie(true);
    }
  }, [cookies]);

  if (!hasCookie) return <></>;

  return (
    <div className='w-full mt-2 bg-gray-300 shadow-md px-1 py-1 rounded-md'>
      <div className='mx-4 my-2'>
        <span>What&apos;s on your mind?</span>
        <form className='flex flex-row mt-2' onSubmit={(e) => handleSubmit(e)}>
          <input
            type='text'
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            className='px-2 py-1 mr-2 rounded-md w-full border shadown-md'
            placeholder='write here...'
          />
          <button
            className='px-3 py-1 border rounded-md shadow-md bg-blue-500 text-white'
            type='submit'>
            Share
          </button>
        </form>
      </div>
    </div>
  );

  async function handleSubmit(e: any) {
    e.preventDefault();

    try {
      if (!hasCookie) return;

      await postTweet({
        message,
        userId: cookies.user_id,
      });
      setMessage("");
      getPublicTweets().then((res) => setTweets(res.tweets));
    } catch (e) {
      console.log("Error: ", e)
    }
  }
}
