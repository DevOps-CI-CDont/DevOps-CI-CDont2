type Methods = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

interface FetchHeadersProps {
	userId?: string;
	method: Methods;
}

export function getFetchHeaders({
	userId,
	method,
}: FetchHeadersProps): HeadersInit {
	return {
		method,

		redirect: "follow",
	};
}
