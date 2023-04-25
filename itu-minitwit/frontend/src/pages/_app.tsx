import useAuthStore from "@/store/authStore";
import useUserStore from "@/store/userStore";
import "@/styles/globals.css";
import { ThemeProvider } from "next-themes";
import type { AppProps } from "next/app";
import { useEffect } from "react";
import { useCookies } from "react-cookie";
import { QueryClient, QueryClientProvider } from "react-query";

export const queryClient = new QueryClient();

export default function App({ Component, pageProps }: AppProps) {
	const [sessionCookie] = useCookies(["user_id"]);

	const setAuthStore = useAuthStore((state) => state.setIsAuth);
	const setUserId = useUserStore((state) => state.setUserId);

	useEffect(() => {
		if (sessionCookie.user_id) {
			setAuthStore(true);
			setUserId(sessionCookie.user_id);
		} else {
			setAuthStore(false);
			setUserId(undefined);
		}
	}, [sessionCookie]);

	return (
		<QueryClientProvider client={queryClient}>
			<ThemeProvider attribute="class">
				<Component {...pageProps} />
			</ThemeProvider>
		</QueryClientProvider>
	);
}
