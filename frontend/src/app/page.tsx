"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

const features = [
  {
    key: "strings",
    label: "String Keys",
    short: "Store, increment, and expire string values.",
    explainer: "Support for SET, GET, DEL, INCR, DECR, EXPIRE, TTL, MSET/MGET and more. Fast in-memory key-value access.",
    commands: ["SET", "GET", "DEL", "MSET", "MGET", "INCR", "DECR", "EXPIRE", "TTL"],
  },
  {
    key: "lists",
    label: "Lists",
    short: "Push, pop, and measure lists efficiently.",
    explainer: "Backed by LPUSH, RPOP, and LLEN. RedisGo allows storing ordered collections, perfect for queues and stacks.",
    commands: ["LPUSH", "RPOP", "LLEN"],
  },
  {
    key: "sets",
    label: "Sets",
    short: "Add and remove unique members.",
    explainer: "SADD, SREM, and SMEMBERS let you build dynamic, duplicate-free groups.",
    commands: ["SADD", "SREM", "SMEMBERS"],
  },
  {
    key: "hashes",
    label: "Hashes",
    short: "Dictionary-style fields for each key.",
    explainer: "Model structured objects with HSET, HGET, HDEL, HGETALL. Popular for user profiles and more.",
    commands: ["HSET", "HGET", "HDEL", "HGETALL"],
  },
  {
    key: "sortedsets",
    label: "Sorted Sets",
    short: "Maintain a ranking with scores.",
    explainer: "Powerful for leaderboards and rankings, using ZADD, ZREM, ZRANGE.",
    commands: ["ZADD", "ZREM", "ZRANGE"],
  },
  {
    key: "introspection",
    label: "Introspection",
    short: "Inspect keys or dump all data.",
    explainer: "KEYS shows all non-expired keys, DUMPALL gets all string keys and values.",
    commands: ["KEYS", "DUMPALL"],
  },
  {
    key: "meta",
    label: "Meta & Health",
    short: "Ping server, echo, get help.",
    explainer: "PING, ECHO, COMMANDS/HELP for developer-friendliness and monitoring.",
    commands: ["PING", "ECHO", "COMMANDS", "HELP"],
  }
];

export default function Home() {
  const [hovered, setHovered] = useState<string | null>(null);
  const router = useRouter();

  return (
    <main className="min-h-screen flex flex-col items-center justify-start px-4 py-16 bg-background text-foreground font-sans">
      <header className="flex flex-col items-center mb-12">
        <Image src="/redis.png" width={64} height={64} alt="RedisGo logo" className="mb-3 rounded shadow-xl" />
        <h1 className="text-4xl font-bold tracking-tighter mb-2">RedisGo</h1>
        <p className="text-lg max-w-xl text-center text-gray-700 dark:text-gray-300 mb-2">
          A minimal, educational Redis clone in Go — fast, in-memory, thread-safe, and ready for you to explore.
        </p>
        <Link
          href="/demo"
          className="mt-3 px-6 py-2 rounded-lg bg-blue-600 text-white text-base font-semibold hover:bg-blue-700 shadow transition"
        >
          Try the Demo &rarr;
        </Link>
      </header>
      <section className="w-full max-w-4xl" id="features">
        <h2 className="text-2xl font-bold mb-7">✨ Features</h2>
        <ul className="grid md:grid-cols-2 gap-5">
          {features.map((f) => (
            <li
              key={f.key}
              className={`relative group bg-[var(--muted)] dark:bg-[var(--muted)] overflow-visible rounded-2xl p-6 border border-transparent hover:border-blue-400 shadow-lg transition cursor-pointer`}
              onMouseEnter={() => setHovered(f.key)}
              onMouseLeave={() => setHovered(null)}
              onClick={() => {
                if (f.key === "strings") {
                  router.push("/demo?feature=strings");
                } else {
                  router.push("/demo");
                }
              }}
              tabIndex={0}
              aria-describedby={`feature-detail-${f.key}`}
            >
              <div className="text-xl font-semibold mb-2 flex items-center gap-2">
                {f.label}
                <span className="inline text-xs text-gray-500 dark:text-gray-300 font-mono">
                  ({f.commands.join(", ")})
                </span>
              </div>
              <div className="text-base text-zinc-600 dark:text-zinc-300">{f.short}</div>
              {/* Expanding explainer on hover/focus */}
              <div
                id={`feature-detail-${f.key}`}
                className={`absolute left-0 right-0 top-full mt-3 z-10 pointer-events-none transition-all duration-300
                  bg-white dark:bg-zinc-900 rounded-xl shadow-xl p-3 text-base text-zinc-900 dark:text-zinc-100
                  ${hovered === f.key ? "opacity-100 translate-y-0 scale-100 visible" : "opacity-0 scale-95 invisible"}
                `}
                style={{ minWidth: "250px" }}
                aria-hidden={hovered !== f.key}
              >
                {f.explainer}
              </div>
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}