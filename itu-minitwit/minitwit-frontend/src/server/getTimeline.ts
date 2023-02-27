  export async function getTimeline(cookie: number) {
    return await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/mytimeline`,
      {
        mode: "no-cors",
        cache: "no-cache",
        method: "GET",
        headers: {
          "Cookie": `user_id=${cookie}`,
          "Content-Type": "application/json",
          origin: "http://localhost:3000",
        },
        credentials: 'same-origin',
        redirect: "follow",
        
      }
    ).then((response) => response.json());
  }

