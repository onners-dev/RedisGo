"use client";
import { useState } from "react";
import axios from "axios";

const API = "http://localhost:8000";

export default function StringDemo() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [response, setResponse] = useState<string | null>(null);
  const [ttl, setTtl] = useState<number | null>(null);

  // Handlers
  const handleSet = async () => {
    await axios.post(`${API}/set`, { key, value });
    setResponse("OK");
  };

  const handleGet = async () => {
    const res = await axios.get<{ value: string | null }>(`${API}/get/${key}`);
    setResponse(res.data.value ?? "<nil>");
  };

  const handleDel = async () => {
    const res = await axios.post(`${API}/cli`, { cmd: `DEL ${key}` });
    setResponse(res.data.resp.trim());
  };

  const handleIncrDecr = async (op: "incr" | "decr") => {
    const res = await axios.post<{ value: number }>(`${API}/counter`, { key, action: op });
    setResponse(res.data.value.toString());
  };

  const handleExpire = async () => {
    const res = await axios.post(`${API}/cli`, { cmd: `EXPIRE ${key} ${ttl ?? 10}` });
    setResponse(`Expire set to ${(ttl ?? 10)}s: ${res.data.resp.trim()}`);
  };

  const handleTtl = async () => {
    const res = await axios.post(`${API}/cli`, { cmd: `TTL ${key}` });
    setResponse(`TTL: ${res.data.resp.trim()}`);
  };

  return (
    <div className="w-full max-w-xl mx-auto flex flex-col gap-4">
      <form
        className="grid grid-cols-1 md:grid-cols-5 gap-2 items-center"
        onSubmit={e => { e.preventDefault(); handleSet(); }}
        autoComplete="off"
      >
        <label className="font-medium text-base text-[var(--foreground)] md:col-span-1">Key</label>
        <input
          className="rounded-lg p-2 border border-[var(--border)] w-full md:col-span-1 text-[var(--foreground)] bg-[var(--background)]"
          placeholder="Key"
          value={key}
          onChange={e => setKey(e.target.value)}
          required
        />
        <label className="font-medium text-base text-[var(--foreground)] md:col-span-1">Value</label>
        <input
          className="rounded-lg p-2 border border-[var(--border)] w-full md:col-span-1 text-[var(--foreground)] bg-[var(--background)]"
          placeholder="Value"
          value={value}
          onChange={e => setValue(e.target.value)}
        />
        <button
          type="submit"
          className="rounded-lg bg-blue-600 text-white px-4 py-2 font-semibold hover:bg-blue-700 w-full md:col-span-1 transition"
        >
          SET
        </button>
      </form>
      <div className="grid grid-cols-2 md:grid-cols-5 gap-2">
        <button
          className="rounded-lg bg-green-600 text-white px-4 py-2 font-semibold hover:bg-green-700 transition"
          onClick={handleGet}
        >
          GET
        </button>
        <button
          className="rounded-lg bg-rose-500 text-white px-4 py-2 font-semibold hover:bg-rose-600 transition"
          onClick={handleDel}
        >
          DEL
        </button>
        <button
          className="rounded-lg bg-gray-300 text-black px-4 py-2 font-semibold hover:bg-gray-400 transition"
          onClick={() => handleIncrDecr("incr")}
        >
          INCR
        </button>
        <button
          className="rounded-lg bg-gray-300 text-black px-4 py-2 font-semibold hover:bg-gray-400 transition"
          onClick={() => handleIncrDecr("decr")}
        >
          DECR
        </button>
        <button
          className="rounded-lg bg-orange-500 text-white px-4 py-2 font-semibold hover:bg-orange-600 transition"
          onClick={handleTtl}
        >
          TTL
        </button>
      </div>
      <div className="flex flex-col md:flex-row items-center gap-3">
        <div className="flex flex-row gap-2 items-center w-full md:w-auto">
          <input
            type="number"
            className="rounded-lg p-2 border border-[var(--border)] w-24 text-[var(--foreground)] bg-[var(--background)]"
            placeholder="Expire (s)"
            value={ttl ?? ""}
            onChange={e => setTtl(Number(e.target.value))}
            min={1}
          />
          <button
            className="rounded-lg bg-violet-600 text-white px-4 py-2 font-semibold hover:bg-violet-700 transition"
            onClick={handleExpire}
          >
            EXPIRE
          </button>
        </div>
        {response !== null && (
          <div className="w-full mt-3 md:mt-0 rounded bg-[var(--muted)] p-3 text-base font-mono text-[var(--foreground)] border border-[var(--border)] shadow transition-all">
            Result: {response}
          </div>
        )}
      </div>
    </div>
  );
}
