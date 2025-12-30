"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Search, Command, Github, Terminal, Cpu, Globe } from "lucide-react";

export default function Home() {
  const [username, setUsername] = useState("");
  const router = useRouter();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      router.push(`/${username}`);
    }
  };

  return (
    <main className="min-h-screen flex flex-col items-center justify-center bg-[#050505] text-white p-4 relative overflow-hidden selection:bg-white/20">
      <div className="fixed inset-0 premium-grid-bg opacity-30 pointer-events-none" />
      <div className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-white/5 blur-[120px] rounded-full pointer-events-none" />

      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
        className="w-full max-w-xl space-y-12 relative z-10"
      >
        <div className="text-center space-y-6">
          <motion.div 
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="inline-flex items-center justify-center w-20 h-20 rounded-3xl bg-gradient-to-b from-white/10 to-white/5 border border-white/10 mb-6 shadow-2xl shadow-black/50 backdrop-blur-xl"
          >
            <Github className="w-10 h-10 text-white" />
          </motion.div>
          
          <motion.h1 
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="text-5xl md:text-7xl font-bold tracking-tighter premium-text-gradient"
          >
            GitHub Stats
          </motion.h1>
          
          <motion.p 
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.4 }}
            className="text-zinc-400 text-xl font-light max-w-md mx-auto leading-relaxed"
          >
            Your GitHub profile, reimagined. <br/>
            Insights wrapped in a stunning visual experience.
          </motion.p>
        </div>

        <motion.form 
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.5 }}
          onSubmit={handleSubmit} 
          className="relative group max-w-md mx-auto"
        >
          <div className="absolute -inset-1 bg-gradient-to-r from-white/20 via-white/10 to-white/20 rounded-xl blur opacity-20 group-hover:opacity-40 transition duration-1000 group-hover:duration-200" />
          <div className="relative flex items-center bg-[#0A0A0A] border border-white/10 rounded-xl p-2 shadow-2xl transition-transform group-hover:scale-[1.02]">
            <div className="pl-4 pr-3 text-zinc-500">
              <Search className="w-5 h-5" />
            </div>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter GitHub username..."
              className="w-full bg-transparent border-none text-white placeholder-zinc-600 focus:ring-0 text-lg h-12 font-medium outline-none"
              autoFocus
            />
            <div className="pr-2">
              <button 
                type="submit"
                disabled={!username.trim()}
                className="px-4 py-2 bg-white text-black rounded-lg font-medium text-sm hover:bg-zinc-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Analyze
              </button>
            </div>
          </div>
        </motion.form>

        <motion.div 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="flex flex-wrap justify-center gap-3 text-sm text-zinc-500"
        >
          <span>Try:</span>
          {['torvalds', 'shadcn', 'vercel'].map((u) => (
            <button
              key={u}
              onClick={() => router.push(`/${u}`)}
              className="hover:text-white transition-colors border-b border-transparent hover:border-white/50 pb-0.5"
            >
              {u}
            </button>
          ))}
        </motion.div>

        <motion.div 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.8, duration: 1 }}
          className="pt-12 border-t border-white/5 flex justify-center gap-12 text-zinc-600"
        >
          <div className="flex items-center gap-2">
            <Terminal className="w-4 h-4" />
            <span className="text-xs font-mono uppercase tracking-widest">Real-time</span>
          </div>
          <div className="flex items-center gap-2">
            <Cpu className="w-4 h-4" />
            <span className="text-xs font-mono uppercase tracking-widest">Fast</span>
          </div>
          <div className="flex items-center gap-2">
            <Globe className="w-4 h-4" />
            <span className="text-xs font-mono uppercase tracking-widest">Global</span>
          </div>
        </motion.div>
      </motion.div>
    </main>
  );
}
