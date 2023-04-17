type Method = "GET";

interface FetchProps {
	url: string;
	method: Method;
	userId?: string;
	formData?: FormData;
}

export async function makeFetch<TData>({
	url,
	method,
	userId,
	formData,
}: FetchProps): Promise<TData> {
	const headers = new Headers();
	headers.append("Content-Type", "application/json");
	userId && headers.append("Authorization", userId);

	const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/${url}`, {
		method,
		headers,
		redirect: "follow",
	});
	return await res.json();
}

export async function makeZodSafeFetch<TData>({
	url,
	method,
}: FetchProps): Promise<TData> {
	return await fetch(`${process.env.NEXT_PUBLIC_API_URL}/${url}`, {
		method,
		headers: {
			"Content-Type": "application/json",
			// Authorization: userId ? `Bearer ${userId}` : "",
		},
		redirect: "follow",
	}).then((res) => res.json());
}
