import Image from "next/image";
import Link from "next/link";

export default function Header() {
  return (
    <header className="w-full py-4 px-6 flex items-center justify-start bg-white dark:bg-zinc-900 border-b dark:border-zinc-800 shadow-sm">
      <Link href="/" className="flex items-center gap-3 group">
        <Image
          src="/redis.png"
          alt="RedisGo logo"
          width={40}
          height={40}
          className="rounded shadow group-hover:scale-105 transition"
        />
        <span className="text-xl font-bold text-zinc-900 dark:text-white tracking-tight group-hover:text-blue-600 transition">
          RedisGo
        </span>
      </Link>
    </header>
  );
}
