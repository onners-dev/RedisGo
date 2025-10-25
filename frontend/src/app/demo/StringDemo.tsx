"use client";
import { useState, useRef, useEffect } from "react";
import axios from "axios";

const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000";
const DEMO_KEY = "foo";
const DEMO_VALUE = "bar";

type KeyPreview = { key: string; value: string | null; ttl: number | null };

function Card({
  icon,
  title,
  description,
  children,
}: {
  icon: string;
  title: string;
  description: string;
  children: React.ReactNode;
}) {
  return (
    <section className="bg-white dark:bg-zinc-900 rounded-2xl border shadow p-6 mb-4">
      <div className="flex items-center gap-2 mb-2">
        <span className="text-xl">{icon}</span>
        <span className="font-bold text-lg">{title}</span>
      </div>
      <p className="text-zinc-500 text-xs mb-4">{description}</p>
      {children}
    </section>
  );
}

export default function StringDemo() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [ttl, setTtl] = useState<number | "">("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Per-action results:
  const [setResult, setSetResult] = useState<string | null>(null);
  const [getResult, setGetResult] = useState<string | null>(null);
  const [delResult, setDelResult] = useState<string | null>(null);
  const [incrResult, setIncrResult] = useState<string | null>(null);
  const [ttlResult, setTtlResult] = useState<string | null>(null);
  const [expireResult, setExpireResult] = useState<string | null>(null);

  // All keys live preview
  const [allKeys, setAllKeys] = useState<KeyPreview[]>([]);
  const timersRef = useRef<{ [k: string]: NodeJS.Timeout }>({});

  // Fetch all keys and values+ttls
  const fetchAllKeys = async () => {
    try {
      // 1. Get all keys
      const keysResp = await axios.post(API + "/cli", { cmd: "KEYS" });
      // The KEYS reply is redis RESP, like: *3\n$1\na\n$1\nb\n$1\nc
      let keys: string[] = [];
      if (typeof keysResp.data.resp === "string" && keysResp.data.resp.startsWith("*")) {
        const lines = keysResp.data.resp.split("\n").map(l => l.trim()).filter(Boolean);
        // extract keys by lines starting with $len then next line is value
        for (let i = 1; i < lines.length; ) {
          if (lines[i].startsWith("$")) {
            keys.push(lines[i + 1]);
            i += 2;
          } else {
            i++;
          }
        }
      }
      // 2. For each, get value and TTL
      const valuesAndTtls = await Promise.all(
        keys.map(async (k) => {
          try {
            const [valueRes, ttlRes] = await Promise.all([
              axios.get<{ value: string | null }>(`${API}/get/${encodeURIComponent(k)}`),
              axios.post(API + "/cli", { cmd: `TTL ${k}` }),
            ]);
            let ttlStr = ttlRes.data.resp?.toString?.().trim();
            let ttl =
              ttlStr && ttlStr.startsWith(":")
                ? parseInt(ttlStr.slice(1))
                : null;
            return { key: k, value: valueRes.data.value, ttl };
          } catch {
            return { key: k, value: null, ttl: null };
          }
        })
      );
      setAllKeys(valuesAndTtls);
    } catch (err: any) {
      setError("Error getting keys.");
      setAllKeys([]);
    }
  };

  // Setup a 1s interval for in-place live updating TTLs
  useEffect(() => {
    // Stop any previous timers
    Object.values(timersRef.current).forEach(clearInterval);
    timersRef.current = {};
    // Setup timer per key with TTL > 0
    setAllKeys((prev) =>
      prev.map((item) =>
        item.ttl && item.ttl > 0 ? { ...item, displayTtl: item.ttl } : item
      )
    );
    if (allKeys.some(k => typeof k.ttl === "number" && k.ttl > 0)) {
      const interval = setInterval(() => {
        setAllKeys((old) =>
          old
            .map((item) => {
              if (typeof item.ttl !== "number" || item.ttl < 0) return item;
              return { ...item, ttl: item.ttl > 0 ? item.ttl - 1 : 0 };
            })
            .filter(i => i.ttl === undefined || i.ttl > 0 || i.ttl === -1)
        );
      }, 1000);
      timersRef.current["__global"] = interval;
    }
    return () => {
      Object.values(timersRef.current).forEach(clearInterval);
      timersRef.current = {};
    };
    // eslint-disable-next-line
  }, [JSON.stringify(allKeys.map((item) => ({ key: item.key, ttl: item.ttl })))]);

  useEffect(() => {
    fetchAllKeys();
    // eslint-disable-next-line
  }, []);

  // Re-fetch all keys after any action
  const handleFetch = () => fetchAllKeys();

  // -- Actions with result, fetch after each
  const run = async (fn: () => Promise<void>, setRes: (msg: string) => void) => {
    setLoading(true);
    setError(null);
    try {
      await fn();
    } catch (err: any) {
      setError(err?.response?.data?.detail || err?.message || "Unknown error");
      setRes("Error occurred!");
    } finally {
      setLoading(false);
      fetchAllKeys();
    }
  };

  return (
    <div className="max-w-lg mx-auto py-8">
      {/* Demo Example */}
      <div className="bg-blue-100 dark:bg-blue-900 text-blue-900 dark:text-blue-100 p-3 rounded text-sm flex flex-wrap gap-2 items-center mb-4">
        <span>Try example:</span>
        <button
          className="underline text-blue-600 dark:text-blue-200 hover:text-blue-900 dark:hover:text-blue-100 transition"
          type="button"
          onClick={() => {
            setKey(DEMO_KEY);
            setValue(DEMO_VALUE);
            setSetResult(null);
            setGetResult(null);
            setDelResult(null);
            setIncrResult(null);
            setTtlResult(null);
            setExpireResult(null);
            setError(null);
            fetchAllKeys();
          }}
        >
          Key: "{DEMO_KEY}", Value: "{DEMO_VALUE}"
        </button>
      </div>

      <Card icon="" title="Set Value" description="Set a value for a given key.">
        <form
          className="flex gap-2 flex-wrap items-center"
          onSubmit={e => {
            e.preventDefault();
            run(async () => {
              await axios.post(`${API}/set`, { key, value });
              setSetResult(`Key "${key}" set to "${value}" 🚀`);
            }, setSetResult);
          }}
        >
          <span className="font-medium">Key</span>
          <input
            className="border p-2 rounded-md flex-1 min-w-[60px] bg-[var(--background)] text-[var(--foreground)]"
            value={key}
            onChange={e => setKey(e.target.value)}
            required
            disabled={loading}
            placeholder="Enter key..."
          />
          <span className="font-medium">Value</span>
          <input
            className="border p-2 rounded-md flex-1 min-w-[60px] bg-[var(--background)] text-[var(--foreground)]"
            value={value}
            onChange={e => setValue(e.target.value)}
            disabled={loading}
            placeholder="Enter value..."
          />
          <button
            type="submit"
            disabled={loading || !key}
            className="rounded-md bg-blue-600 text-white px-4 py-2 font-semibold hover:bg-blue-700 focus:ring-2 focus:ring-blue-400 transition disabled:opacity-60"
          >
            SET
          </button>
        </form>
        {setResult && (
          <div className="mt-2 text-blue-800 dark:text-blue-200 bg-blue-50 dark:bg-blue-900 font-mono rounded p-2">
            {setResult}
          </div>
        )}
      </Card>

      <Card icon="" title="Get Value" description="Fetch the value for a given key.">
        <div className="flex gap-2 flex-wrap items-center">
          <input
            className="border p-2 rounded-md flex-1 min-w-[110px] bg-[var(--background)] text-[var(--foreground)]"
            value={key}
            onChange={e => setKey(e.target.value)}
            placeholder="Key"
            disabled={loading}
          />
          <button
            type="button"
            onClick={() =>
              run(async () => {
                const res = await axios.get<{ value: string | null }>(`${API}/get/${key}`);
                if (res.data.value !== null)
                  setGetResult(`Value for "${key}": "${res.data.value}"`);
                else setGetResult(`Key "${key}" not found`);
              }, setGetResult)
            }
            className="rounded-md bg-green-600 text-white px-4 py-2 font-semibold hover:bg-green-700 transition disabled:opacity-60"
            disabled={loading || !key}
          >
            GET
          </button>
        </div>
        {getResult && (
          <div className="mt-2 text-green-800 dark:text-green-200 bg-green-50 dark:bg-green-900 font-mono rounded p-2">
            {getResult}
          </div>
        )}
      </Card>

      <Card icon="" title="Delete Key" description="Remove a key and its value.">
        <div className="flex gap-2 flex-wrap items-center">
          <input
            className="border p-2 rounded-md flex-1 min-w-[110px]  bg-[var(--background)] text-[var(--foreground)]"
            value={key}
            onChange={e => setKey(e.target.value)}
            placeholder="Key"
            disabled={loading}
          />
          <button
            type="button"
            onClick={() =>
              run(async () => {
                const res = await axios.post(`${API}/cli`, { cmd: `DEL ${key}` });
                setDelResult(
                  res.data.resp.trim() === ":1"
                    ? `Key "${key}" deleted 🗑️`
                    : `Key "${key}" was not found or already deleted`
                );
              }, setDelResult)
            }
            className="rounded-md bg-red-500 text-white px-4 py-2 font-semibold hover:bg-red-600 transition disabled:opacity-60"
            disabled={loading || !key}
          >
            DEL
          </button>
        </div>
        {delResult && (
          <div className="mt-2 text-rose-800 dark:text-rose-200 bg-rose-50 dark:bg-rose-900 font-mono rounded p-2">
            {delResult}
          </div>
        )}
      </Card>

      <Card icon="" title="Increment / Decrement" description="Increase or decrease integer value for a key. The value of the key has to be an integer">
        <div className="flex gap-2 flex-wrap items-center">
          <input
            className="border p-2 rounded-md flex-1 min-w-[110px]  bg-[var(--background)] text-[var(--foreground)]"
            value={key}
            onChange={e => setKey(e.target.value)}
            placeholder="Key"
            disabled={loading}
          />
          <button
            type="button"
            onClick={() =>
              run(async () => {
                const res = await axios.post<{ value: number }>(`${API}/counter`, { key, action: "incr" });
                setIncrResult(`INCR "${key}": ${res.data.value}`);
              }, setIncrResult)
            }
            className="rounded-md bg-gray-700 text-white px-4 py-2 font-semibold hover:bg-gray-900 transition disabled:opacity-60"
            disabled={loading || !key}
          >
            INCR
          </button>
          <button
            type="button"
            onClick={() =>
              run(async () => {
                const res = await axios.post<{ value: number }>(`${API}/counter`, { key, action: "decr" });
                setIncrResult(`DECR "${key}": ${res.data.value}`);
              }, setIncrResult)
            }
            className="rounded-md bg-gray-300 text-black px-4 py-2 font-semibold hover:bg-gray-400 transition disabled:opacity-60"
            disabled={loading || !key}
          >
            DECR
          </button>
        </div>
        {incrResult && (
          <div className="mt-2 text-gray-700 dark:text-gray-100 bg-gray-50 dark:bg-zinc-800 font-mono rounded p-2">
            {incrResult}
          </div>
        )}
      </Card>

      <Card icon="" title="Get TTL" description="Check remaining time to live for a key.">
        <div className="flex gap-2 flex-wrap items-center">
          <input
            className="border p-2 rounded-md flex-1 min-w-[110px]  bg-[var(--background)] text-[var(--foreground)]"
            value={key}
            onChange={e => setKey(e.target.value)}
            placeholder="Key"
            disabled={loading}
          />
          <button
            type="button"
            onClick={() =>
              run(async () => {
                const res = await axios.post(`${API}/cli`, { cmd: `TTL ${key}` });
                const out = res.data.resp.trim();
                setTtlResult(
                  out.startsWith(":")
                    ? `TTL for "${key}": ${out.replace(":", "")} seconds`
                    : `TTL unknown`
                );
              }, setTtlResult)
            }
            className="rounded-md bg-orange-500 text-white px-4 py-2 font-semibold hover:bg-orange-600 transition disabled:opacity-60"
            disabled={loading || !key}
          >
            TTL
          </button>
        </div>
        {ttlResult && (
          <div className="mt-2 text-orange-700 dark:text-orange-200 bg-orange-50 dark:bg-orange-900 font-mono rounded p-2">
            {ttlResult}
          </div>
        )}
      </Card>

      <Card icon="" title="Set Expiry" description="Set a time-to-live (expiry, in seconds) on a key.">
        <form
          className="flex flex-wrap gap-2 items-center"
          onSubmit={e => {
            e.preventDefault();
            run(async () => {
              const seconds = ttl || 10;
              const res = await axios.post(`${API}/cli`, { cmd: `EXPIRE ${key} ${seconds}` });
              setExpireResult(
                res.data.resp.trim() === ":1"
                  ? `Expiry set for "${key}" (${seconds}s) ⏰`
                  : `Couldn't set expiry! Does the key exist?`
              );
            }, setExpireResult);
          }}
        >
          <input
            className="border p-2 rounded-md flex-1 min-w-[70px] bg-[var(--background)] text-[var(--foreground)]"
            value={key}
            onChange={e => setKey(e.target.value)}
            placeholder="Key"
            disabled={loading}
          />
          <input
            type="number"
            min={1}
            step={1}
            className="border p-2 rounded-md w-24 bg-[var(--background)] text-[var(--foreground)]"
            value={ttl}
            onChange={e => setTtl(Number(e.target.value))}
            disabled={loading}
            placeholder="Seconds"
          />
          <button
            type="submit"
            className="rounded-md bg-violet-600 text-white px-4 py-2 font-semibold hover:bg-violet-700 transition disabled:opacity-60"
            disabled={loading || !key}
          >
            EXPIRE
          </button>
        </form>
        {expireResult && (
          <div className="mt-2 text-violet-700 dark:text-violet-200 bg-violet-50 dark:bg-violet-900 font-mono rounded p-2">
            {expireResult}
          </div>
        )}
      </Card>

      {/* Live Preview for ALL KEYS */}
      <div className="border-2 border-blue-400 rounded bg-blue-50 dark:bg-zinc-900 p-4 my-4 font-mono text-blue-800 dark:text-blue-200 shadow">
        <div className="flex justify-between items-center mb-2">
          <span className="text-xs font-semibold text-blue-800 dark:text-blue-300">
            Live Preview: All Current Keys
          </span>
          <button
            onClick={handleFetch}
            className="px-3 py-1 ml-1 rounded bg-blue-200 dark:bg-blue-800 text-blue-900 dark:text-blue-100 text-xs shadow font-semibold hover:bg-blue-300 hover:dark:bg-blue-700 transition"
          >
            Refresh
          </button>
        </div>
        {allKeys.length === 0 ? (
          <div className="italic text-zinc-400 text-sm">No keys in memory.</div>
        ) : (
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="border-b border-blue-200 dark:border-zinc-700">
                <th className="py-1">Key</th>
                <th className="py-1">Value</th>
                <th className="py-1">TTL</th>
              </tr>
            </thead>
            <tbody>
              {allKeys.map((item) => (
                <tr key={item.key} className="border-b border-blue-100 dark:border-zinc-800">
                  <td className="py-1 pr-2 break-all font-bold">{`"${item.key}"`}</td>
                  <td className="py-1 pr-2">{item.value === null ? <span className="italic text-zinc-400">N/A</span> : `"${item.value}"`}</td>
                  <td className="py-1">
                    {typeof item.ttl === "number"
                      ? item.ttl > 0
                        ? <span>{item.ttl}s</span>
                        : <span className="text-rose-400">Expired</span>
                      : <span className="italic text-zinc-400">∞</span>
                    }
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {error && (
        <div className="mt-2 rounded bg-rose-50 text-rose-800 p-3 border border-rose-300 font-mono">{error}</div>
      )}
    </div>
  );
}
