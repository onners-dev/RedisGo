"use client";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import StringDemo from "./StringDemo";



function ListDemo() {
  return <div>📋 <b>List Demo:</b> LPUSH, RPOP, LLEN, etc.</div>;
}
function SetDemo() {
  return <div>🔗 <b>Set Demo:</b> SADD, SREM, SMEMBERS, etc.</div>;
}
function HashDemo() {
  return <div>♻️ <b>Hash Demo:</b> HSET, HGET, HDEL, HGETALL, etc.</div>;
}
function ZSetDemo() {
  return <div>🏅 <b>ZSet Demo:</b> ZADD, ZREM, ZRANGE, etc.</div>;
}
function TTLDemo() {
  return <div>⏳ <b>TTL Demo:</b> KEYS, EXPIRE, TTL, etc.</div>;
}

const features = [
  { key: "strings", label: "String Keys", desc: "SET, GET, DEL, INCR, TTL", Demo: StringDemo },
  { key: "lists", label: "Lists", desc: "LPUSH, RPOP, LLEN", Demo: ListDemo },
  { key: "sets", label: "Sets", desc: "SADD, SREM, SMEMBERS", Demo: SetDemo },
  { key: "hashes", label: "Hashes", desc: "HSET, HGET, HDEL, HGETALL", Demo: HashDemo },
  { key: "zsets", label: "Sorted Sets", desc: "ZADD, ZREM, ZRANGE", Demo: ZSetDemo },
  { key: "ttl", label: "Expiry & TTL", desc: "EXPIRE, TTL", Demo: TTLDemo },
];


export default function VisualDemo() {
  const searchParams = useSearchParams();
  const [selected, setSelected] = useState("strings");
  const CurrentDemo = features.find(f => f.key === selected)?.Demo ?? StringDemo;

  useEffect(() => {
    const feature = searchParams.get("feature");
    if (feature) setSelected(feature);
  }, [searchParams]);

  return (
    <div className="bg-white dark:bg-zinc-900/70 rounded-2xl shadow-lg flex flex-col md:flex-row overflow-hidden min-h-[350px]">
      {/* Sidebar */}
      <nav className="w-full md:w-56 bg-[var(--muted)] dark:bg-zinc-800 flex flex-row md:flex-col p-3 gap-2 border-r dark:border-zinc-700 shrink-0">
        {features.map((f) => (
          <button
            key={f.key}
            onClick={() => setSelected(f.key)}
            className={`w-full text-left px-3 py-2 rounded ${
              selected === f.key
                ? "bg-blue-600 text-white font-semibold shadow"
                : "hover:bg-blue-100 dark:hover:bg-zinc-700"
            } transition`}
            aria-current={selected === f.key}
          >
            <span className="text-base">{f.label}</span>
            <div className="text-xs text-zinc-600 dark:text-zinc-300">{f.desc}</div>
          </button>
        ))}
      </nav>
      {/* Demo area */}
      <section className="flex-1 p-6 md:p-10 flex flex-col gap-5 justify-start">
        <CurrentDemo />
      </section>
    </div>
  );
}
