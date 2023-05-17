import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export async function middleware(req: NextRequest) {
	const cookie = req.cookies.get("user_id");

	if (!cookie) {
		const url = req.nextUrl.clone();
		url.pathname = "/login";
		return NextResponse.redirect(url);
	}
}

export const config = {
	matcher: ["/"],
};
