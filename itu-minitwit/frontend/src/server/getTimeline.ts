interface GetIsFollowingProps {
  userId: string;
}

export async function getTimeline({ userId }: GetIsFollowingProps) {
  var myHeaders = new Headers();
  myHeaders.append("Cookie", `user_id=${userId}`);

  return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/mytimeline`, {
    method: "GET",
    headers: myHeaders,
    redirect: "follow",
  }).then((response) => response.json());
}
