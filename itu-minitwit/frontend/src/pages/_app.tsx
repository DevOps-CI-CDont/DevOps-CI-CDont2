import useAuthStore from "@/store/authStore";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { useEffect } from "react";
import { useCookies } from "react-cookie";

export default function App({ Component, pageProps }: AppProps) {
  const [sessionCookie] = useCookies(["user_id"]);

  const setAuthStore = useAuthStore((state) => state.setIsAuth);

  useEffect(() => {
    if (sessionCookie.user_id) {
      setAuthStore(true);
    } else {
      setAuthStore(false);
    }
  }, [sessionCookie]);

  return <Component {...pageProps} />;
}
