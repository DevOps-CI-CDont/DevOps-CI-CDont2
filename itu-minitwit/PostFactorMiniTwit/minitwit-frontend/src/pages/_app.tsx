import useAuthStore from "@/store/authStore";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { useEffect } from "react";
import { useCookies } from "react-cookie";

export default function App({ Component, pageProps }: AppProps) {
  const [sessionCookie] = useCookies(["session"]);

  const setLoginStore = useAuthStore((state) => state.setIsAuth);

  useEffect(() => {
    if (sessionCookie.session) {
      setLoginStore(true);
    } else {
      setLoginStore(false);
    }
  }, [sessionCookie]);

  return <Component {...pageProps} />;
}
