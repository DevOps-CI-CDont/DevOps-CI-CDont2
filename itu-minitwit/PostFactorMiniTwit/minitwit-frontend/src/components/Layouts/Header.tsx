import { router } from "@/globals/router";
import Link from "next/link";

export function Header() {
  return (
    <div className='border-b shadow-md w-screen fixed top-0 left-0 right-0'>
      <nav className='flex justify-between items-center h-20 max-w-7xl mx-auto px-2'>
        <h2 className='font-bold text-lg'>ITU Minitwit</h2>
        <ul className='flex justify-center items-center'>
          {router.map((route) => (
            <li key={route.id} className='mx-2 hover:underline'>
              <Link href={route.path || ""}>{route.text}</Link>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  );
}
