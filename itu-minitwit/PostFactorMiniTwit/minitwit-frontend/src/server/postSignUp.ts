interface postSignUpProps {
  username: string;
  password: string;
  password2: string;
  email: string;
}

export async function postSignUp({
  username,
  password,
  password2,
  email,
}: postSignUpProps) {
  var formdata = new FormData();
  formdata.append("username", username);
  formdata.append("email", email);
  formdata.append("password", password);
  formdata.append("password2", password2);

  return await fetch(
    `${process.env.NEXT_PUBLIC_CORS_ORIGIN}/${process.env.NEXT_PUBLIC_API_URL}/register`,
    {
      method: "POST",
      redirect: "follow",
      body: formdata,
    }
  ).then((response) => response.json());
}
