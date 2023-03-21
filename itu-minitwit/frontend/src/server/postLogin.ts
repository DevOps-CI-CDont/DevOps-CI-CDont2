interface PostLoginProps {
	username: string;
	password: string;
}

export async function postLogin({ username, password }: PostLoginProps) {
	var formdata = new FormData();
	formdata.append("username", username);
	formdata.append("password", password);

	const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/login`, {
		method: "POST",
		body: formdata,
		redirect: "follow",
	}).then((response) => response.json());

	return res;
}
