interface PostLoginProps {
  username: string;
  password: string;
}

export async function postLogin(username: string, password: string) {
  var myHeaders = new Headers();
  //myHeaders.append("Cookie", "user_id=3");

  var formdata = new FormData();
  formdata.append("username", username);
  formdata.append("password", password);

  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/login`,
    {
      method: "POST",
      mode: "no-cors",
      headers: myHeaders,
      body: formdata,
    }
  ).then((res) => res.json());

  return res;
}
