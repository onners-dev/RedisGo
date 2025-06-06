"use client";

import { useState } from "react";
import axios from "axios";
import Image from "next/image";

export default function Home() {
  // States for set/get
  const [setKey, setSetKey] = useState("");
  const [setValue, setSetValue] = useState("");
  const [getKey, setGetKey] = useState("");
  const [getValue, setGetValue] = useState<string | null>(null);

  // States for counter
  const [counter, setCounter] = useState<number | null>(null);

  // States for keys
  const [keys, setKeys] = useState<string[]>([]);
  const [loadingKeys, setLoadingKeys] = useState(false);

  // Replace with your Python backend URL
  const API = "http://localhost:8000";

  // Handlers
  const handleSet = async () => {
    if (!setKey) return;
    await axios.post(`${API}/set`, { key: setKey, value: setValue });
    setSetKey("");
    setSetValue("");
  };
  const handleGet = async () => {
    if (!getKey) return;
    const res = await axios.get<{ value: string | null }>(`${API}/get/${getKey}`);
    setGetValue(res.data.value);
  };

  // Counter handlers (assumes you implement this API in Python)
  const handleUpdateCounter = async (action: "incr" | "decr") => {
    const key = "counter";
    const res = await axios.post<{ value: number }>(
      `${API}/counter`,
      { key, action }
    );
    setCounter(res.data.value);
  };

  // Load keys (assumes you implement this API in Python)
  const fetchKeys = async () => {
    setLoadingKeys(true);
    const res = await axios.get<{ keys: string[] }>(`${API}/keys`);
    setKeys(res.data.keys);
    setLoadingKeys(false);
  };

  return (
    <div className="min-h-screen bg-background text-foreground px-4 py-14 flex flex-col justify-start items-center font-sans">
      <header className="mb-10 flex items-center gap-4">
        <Image
          src="/redis.png"
          alt="RedisGo Logo"
          width={48}
          height={48}
          className="rounded shadow"
        />
        <h1 className="text-3xl font-bold tracking-tight">RedisGo Showcase</h1>
      </header>

      <div className="grid gap-7 w-full max-w-2xl">
        {/* Set/Get Section */}
        <section className="bg-white/70 dark:bg-zinc-900/80 rounded-xl shadow p-6 flex flex-col gap-4">
          <h2 className="text-xl font-semibold mb-1">📝 Set / Get String</h2>
          <div className="flex gap-2 items-end">
            <div>
              <label className="text-xs">Set key</label>
              <input
                type="text"
                value={setKey}
                onChange={e => setSetKey(e.target.value)}
                className="rounded p-2 border w-28 mr-2 text-black"
                placeholder="key"
              />
              <input
                type="text"
                value={setValue}
                onChange={e => setSetValue(e.target.value)}
                className="rounded p-2 border w-28 mr-2 text-black"
                placeholder="value"
              />
              <button
                className="rounded bg-blue-500 text-white px-3 py-2 hover:bg-blue-600 transition"
                onClick={handleSet}
              >
                SET
              </button>
            </div>
            <div>
              <label className="text-xs">Get key</label>
              <input
                type="text"
                value={getKey}
                onChange={e => setGetKey(e.target.value)}
                className="rounded p-2 border w-28 mr-2 text-black"
                placeholder="key"
              />
              <button
                className="rounded bg-green-600 text-white px-3 py-2 hover:bg-green-700 transition"
                onClick={handleGet}
              >
                GET
              </button>
              <span className="ml-2 text-sm font-mono">
                {getValue !== null && <span>Value: <b>{getValue}</b></span>}
              </span>
            </div>
          </div>
        </section>

        {/* Counter */}
        <section className="bg-white/70 dark:bg-zinc-900/80 rounded-xl shadow p-6 flex flex-col gap-3">
          <h2 className="text-xl font-semibold mb-2">🔢 INCR / DECR (Counter Demo)</h2>
          <div className="flex gap-4 items-center">
            <button
              className="rounded bg-gray-200 dark:bg-zinc-800 px-4 py-2 text-2xl font-bold"
              onClick={() => handleUpdateCounter("decr")}
            >-</button>
            <span className="text-2xl">{counter !== null ? counter : "--"}</span>
            <button
              className="rounded bg-gray-200 dark:bg-zinc-800 px-4 py-2 text-2xl font-bold"
              onClick={() => handleUpdateCounter("incr")}
            >+</button>
            <button
              className="ml-4 text-sm underline"
              onClick={async () => {
                const res = await axios.get<{ value: string | null }>(`${API}/get/counter`);
                setCounter(res.data.value ? parseInt(res.data.value) : 0);
              }}
            >
              Load value
            </button>
          </div>
        </section>

        {/* Keys */}
        <section className="bg-white/70 dark:bg-zinc-900/80 rounded-xl shadow p-6">
          <h2 className="text-xl font-semibold mb-2">🔑 All Keys</h2>
          <button
            className="rounded bg-violet-500 text-white px-3 py-2 hover:bg-violet-600 transition mb-3"
            onClick={fetchKeys}
            disabled={loadingKeys}
          >
            {loadingKeys ? "Loading..." : "Show Keys"}
          </button>
          <div className="flex flex-wrap gap-2 mt-2">
            {keys.length === 0 && !loadingKeys ? (
              <span className="text-zinc-500 text-sm">No keys.</span>
            ) : (
              keys.map((k, i) => (
                <span key={i} className="rounded bg-zinc-200 dark:bg-zinc-800 px-2 py-1 font-mono text-xs">{k}</span>
              ))
            )}
          </div>
        </section>
      </div>

      <footer className="mt-16 text-sm text-zinc-500 text-center">
        Powered by <a href="https://github.com/yourusername/RedisGo" className="underline hover:text-blue-700">RedisGo</a>
      </footer>
    </div>
  );
}
