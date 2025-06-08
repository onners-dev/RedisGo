"use client";
import { useState, useRef, useEffect } from "react";
import axios from "axios";

const API = "http://localhost:8000";

export default function CommandLineDemo() {
  const [history, setHistory] = useState<string[]>(["Welcome to RedisGo! Type a command:"]);
  const [input, setInput] = useState("");
  const endRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [history]);

  async function handleCommand(e: React.FormEvent) {
    e.preventDefault();
    if (!input.trim()) return;
    setHistory(h => [...h, `> ${input}`]);
    try {
      const res = await axios.post(`${API}/cli`, { cmd: input });
      setHistory(h => [...h, res.data.resp || "<no response>"]);
    } catch (err: any) {
      setHistory(h => [...h, err.response?.data?.detail || "Network error"]);
    }
    setInput("");
  }

  return (
    <div className="bg-black text-green-300 font-mono rounded-xl p-5 shadow-lg min-h-[350px] max-h-[450px] overflow-y-auto">
      {history.map((line, i) => (
        <div key={i} className="whitespace-pre-wrap break-words">{line}</div>
      ))}
      <form className="pt-2 flex gap-2" onSubmit={handleCommand}>
        <span className="text-green-400 font-bold">{">"}</span>
        <input
          className="bg-black text-green-200 border-0 focus:outline-none font-mono flex-1"
          autoFocus
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder="Type command and hit Enter..."
        />
      </form>
      <div ref={endRef} />
    </div>
  );
}
