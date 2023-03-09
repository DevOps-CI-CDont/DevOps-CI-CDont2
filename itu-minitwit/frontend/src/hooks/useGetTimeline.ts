import { getTimeline } from "@/server/getTimeline";
import { useQuery } from "react-query";

export function useGetTimeline(userId: string) {
  const { isLoading, error, data } = useQuery(
    ["timeline", userId],
    async () => {
      var myHeaders = new Headers();
      myHeaders.append("Cookie", `user_id=${userId}`);

      return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/mytimeline`, {
        method: "GET",
        headers: myHeaders,
        redirect: "follow",
      }).then((response) => response.text());
    }
  );

  return {
    isLoading,
    error,
    data,
  };
}
