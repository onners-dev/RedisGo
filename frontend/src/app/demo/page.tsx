"use client";
import { useState } from "react";
import CommandLineDemo from "./CommandLineDemo";
import VisualDemo from "./VisualDemo";

export default function DemoPage() {
  const [mode, setMode] = useState<"cli" | "visual">("cli");

  return (
    <main className="min-h-screen px-4 py-12 bg-background text-foreground flex flex-col items-center">
      <h1 className="text-3xl font-bold mb-6">RedisGo</h1>
      <div className="mb-6">
        <button
          className={`px-5 py-2 rounded-l-lg ${mode === "cli"
            ? "bg-blue-600 text-white"
            : "bg-zinc-200 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300"
          } font-semibold transition`}
          onClick={() => setMode("cli")}
        >
          Command Line
        </button>
        <button
          className={`px-5 py-2 rounded-r-lg ${mode === "visual"
            ? "bg-blue-600 text-white"
            : "bg-zinc-200 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300"
          } font-semibold transition`}
          onClick={() => setMode("visual")}
        >
          Visual Playground
        </button>
      </div>
      <div className="w-full max-w-2xl">
        {mode === "cli" ? <CommandLineDemo /> : <VisualDemo />}
      </div>
    </main>
  );
}
