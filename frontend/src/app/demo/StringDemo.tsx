"use client";
import { useState } from "react";
import axios from "axios";

const API = "http://localhost:8000";
const DEMO_KEY = "foo", DEMO_VALUE = "bar";

export default function StringDemo() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [response, setResponse] = useState<string | null>(null);
  const [ttl, setTtl] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);

  // Helper to run async actions with loading UI and error handling
  const run = async (fn: () => Promise<any>) => {
    setLoading(true);
    try {
      await fn();
    } finally {
      setLoading(false);
    }
  };

  // Handlers
  const handleSet = async () => {
    await axios.post(`${API}/set`, { key, value });
    setResponse(`Key "${key}" set to "${value}" 🚀`);
  };

  const handleGet = async () => {
    const res = await axios.get<{ value: string | null }>(`${API}/get/${key}`);
    setResponse(
      res.data.value !== null
        ? `Value for "${key}": "${res.data.value}"`
        : `Key "${key}" not found`
    );
  };

  const handleDel = async () => {
    const res = await axios.post(`${API}/cli`, { cmd: `DEL ${key}` });
    const out = res.data.resp.trim();
    setResponse(out === ":1" ? `Key "${key}" deleted 🗑️` : `Key "${key}" was not found or already deleted`);
  };

  const handleIncrDecr = async (op: "incr" | "decr") => {
    const res = await axios.post<{ value: number }>(`${API}/counter`, { key, action: op });
    setResponse(`${op === "incr" ? "INCR" : "DECR"} "${key}": ${res.data.value}`);
  };

  const handleExpire = async () => {
    const seconds = ttl ?? 10;
    const res = await axios.post(`${API}/cli`, { cmd: `EXPIRE ${key} ${seconds}` });
    setResponse(
      res.data.resp.trim() === ":1"
        ? `Expiry set for "${key}" (${seconds} seconds) ⏰`
        : `Couldn't set expiry! Does the key exist?`
    );
  };

  const handleTtl = async () => {
    const res = await axios.post(`${API}/cli`, { cmd: `TTL ${key}` });
    const out = res.data.resp.trim();
    setResponse(
      out.startsWith(":")
        ? `TTL for "${key}": ${out.replace(":", "")} seconds`
        : `TTL unknown`
    );
  };

  return (
    <div className="w-full max-w-xl mx-auto flex flex-col gap-6 py-6">
      {/* Demo Example */}
      <div className="bg-blue-100 dark:bg-blue-900 text-blue-900 dark:text-blue-100 p-3 rounded text-sm flex flex-wrap gap-2 items-center">
        <span>Try a demo:&nbsp;</span>
        <button
          className="underline text-blue-600 dark:text-blue-300 hover:text-blue-900 dark:hover:text-blue-100 transition"
          type="button"
          onClick={() => {
            setKey(DEMO_KEY);
            setValue(DEMO_VALUE);
            setResponse(null);
          }}
        >
          Key: "{DEMO_KEY}", Value: "{DEMO_VALUE}"
        </button>
      </div>

      {/* SET Section */}
      <form
        onSubmit={e => {
          e.preventDefault();
          run(handleSet);
        }}
        className="flex flex-wrap gap-2 items-center p-4 bg-[var(--muted)] rounded shadow border border-[var(--border)]"
        autoComplete="off"
      >
        <label className="font-medium text-base text-[var(--foreground)]">Key</label>
        <input
          className="p-2 rounded-md border border-[var(--border)] flex-1 min-w-[120px] text-[var(--foreground)] bg-[var(--background)]"
          value={key}
          onChange={e => setKey(e.target.value)}
          required
          disabled={loading}
          placeholder="Enter key..."
        />
        <label className="font-medium text-base text-[var(--foreground)]">Value</label>
        <input
          className="p-2 rounded-md border border-[var(--border)] flex-1 min-w-[120px] text-[var(--foreground)] bg-[var(--background)]"
          value={value}
          onChange={e => setValue(e.target.value)}
          disabled={loading}
          placeholder="Enter value..."
        />
        <button
          type="submit"
          disabled={loading}
          className="rounded-md bg-blue-600 text-white px-4 py-2 font-semibold hover:bg-blue-700 focus:ring-2 focus:ring-blue-400 transition disabled:opacity-60"
        >
          SET
        </button>
      </form>

      {/* Action Buttons */}
      <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
        <button
          className="rounded-md bg-green-600 text-white px-4 py-2 font-semibold hover:bg-green-700 transition disabled:opacity-60"
          disabled={loading || !key}
          onClick={() => run(handleGet)}
        >
          GET
        </button>
        <button
          className="rounded-md bg-rose-500 text-white px-4 py-2 font-semibold hover:bg-rose-600 transition disabled:opacity-60"
          disabled={loading || !key}
          onClick={() => run(handleDel)}
        >
          DEL
        </button>
        <button
          className="rounded-md bg-gray-300 text-black px-4 py-2 font-semibold hover:bg-gray-400 transition disabled:opacity-60"
          disabled={loading || !key}
          onClick={() => run(() => handleIncrDecr("incr"))}
        >
          INCR
        </button>
        <button
          className="rounded-md bg-gray-300 text-black px-4 py-2 font-semibold hover:bg-gray-400 transition disabled:opacity-60"
          disabled={loading || !key}
          onClick={() => run(() => handleIncrDecr("decr"))}
        >
          DECR
        </button>
        <button
          className="rounded-md bg-orange-500 text-white px-4 py-2 font-semibold hover:bg-orange-600 transition disabled:opacity-60"
          disabled={loading || !key}
          onClick={() => run(handleTtl)}
        >
          TTL
        </button>
      </div>

      {/* Expiry Section */}
      <div className="flex flex-col md:flex-row gap-3 items-center bg-[var(--muted)] rounded p-4 border border-[var(--border)]">
        <label className="font-medium text-base text-[var(--foreground)]">Set Expiry (sec):</label>
        <input
          type="number"
          min={1}
          className="p-2 rounded-md border border-[var(--border)] w-28 text-[var(--foreground)] bg-[var(--background)]"
          placeholder="Expire (s)"
          value={ttl ?? ""}
          onChange={e => setTtl(Number(e.target.value))}
          disabled={loading}
        />
        <button
          className="rounded-md bg-violet-600 text-white px-4 py-2 font-semibold hover:bg-violet-700 transition disabled:opacity-60"
          disabled={loading || !key}
          onClick={() => run(handleExpire)}
          type="button"
        >
          EXPIRE
        </button>
      </div>

      {/* Result */}
      {loading ? (
        <div className="text-center text-lg text-gray-400 p-6 select-none">Loading...</div>
      ) : response !== null && (
        <div className="rounded bg-black text-green-300 font-mono text-lg p-4 mt-2 shadow border-2 border-green-500 select-all">
          <b>Result:</b> {response}
        </div>
      )}
    </div>
  );
}
