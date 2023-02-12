import { Footer } from "@/components/Layouts/Footer";
import { Header } from "@/components/Layouts/Header";

interface DefaultLayoutProps {
  children?: React.ReactNode;
}

export default function DefaultLayout({ children }: DefaultLayoutProps) {
  return (
    <div>
      <Header />
      <main className='pt-20 pb-20 mx-auto max-w-3xl min-h-screen px-2'>
        {children}
      </main>
      <Footer />
    </div>
  );
}
