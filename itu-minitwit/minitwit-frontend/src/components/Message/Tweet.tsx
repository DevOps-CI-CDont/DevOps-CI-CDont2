import { getAvatarUrl } from "@/lib/dicebear";
import { Tweet } from "@/types/tweet.type";
import Image from "next/image";
import Link from "next/link";

export function Tweet({ text, author, pub_date }: Tweet) {
  return (
    <div className='mx-4 my-4 py-2 px-2 bg-blue-500 flex rounded-md'>
      <div>
        <Image
          src={getAvatarUrl(author.username)}
          alt='avatar'
          width={60}
          height={60}
          className='rounded-full shadow-md'
        />
      </div>
      <div className='ml-4 my-2 text-white flex flex-col'>
        <span className='text-md'>{text}</span>
        <span className='text-sm text-gray-200'>
          Tweeted by{" "}
          <Link href={`/user/${author.username}`} className='font-bold'>
            {author.username}
          </Link>{" "}
          at {new Date(pub_date).toDateString()}
        </span>
      </div>
    </div>
  );
}
