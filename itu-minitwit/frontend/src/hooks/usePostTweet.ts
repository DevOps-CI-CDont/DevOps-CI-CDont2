import { useMutation } from "react-query";
import { postTweet, PostTweetProps } from "@/server/postTweet";
import { queryClient } from "@/pages/_app";

export async function usePostTweet({ message, userId }: PostTweetProps) {
  const mutation = useMutation(await postTweet({ message, userId }), {
    onSuccess: () => {
      queryClient.invalidateQueries("timeline");
    },
  });

  return mutation;
}
