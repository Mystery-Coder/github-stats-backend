"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { motion } from "framer-motion";
import { 
  MapPin, Link as LinkIcon, Twitter, Star, GitFork, 
  BookOpen, Box, Activity, Layers, Users, Calendar, 
  ArrowUpRight, Github 
} from "lucide-react";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";
import { UserStats, PinnedRepos } from "../types";

const container = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1
    }
  }
};

const item = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0 }
};

export default function DashboardPage() {
  const params = useParams();
  const username = params.username as string;
  const [stats, setStats] = useState<UserStats | null>(null);
  const [pinned, setPinned] = useState<PinnedRepos | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [hoveredLang, setHoveredLang] = useState<any>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [statsRes, pinnedRes] = await Promise.all([
          fetch(`http://localhost:8080/api/stats/${username}`),
          fetch(`http://localhost:8080/api/pinned/${username}`)
        ]);

        if (!statsRes.ok || !pinnedRes.ok) throw new Error("Failed to fetch data");

        const statsData = await statsRes.json();
        const pinnedData = await pinnedRes.json();

        if (statsData.errors) throw new Error(statsData.errors[0].message);
        
        setStats(statsData);
        setPinned(pinnedData);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (username) fetchData();
  }, [username]);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#050505] flex items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <div className="w-12 h-12 border-[3px] border-white/10 border-t-white rounded-full animate-spin" />
          <p className="text-zinc-500 text-sm font-mono animate-pulse">GENERATING INSIGHTS...</p>
        </div>
      </div>
    );
  }

  if (error || !stats || !pinned) {
    return (
      <div className="min-h-screen bg-[#050505] flex items-center justify-center text-red-400 font-mono">
        ERROR: {error || "USER_NOT_FOUND"}
      </div>
    );
  }

  const user = stats.user;
  
  // Calculate language stats
  const languageMap = new Map<string, { size: number; color: string }>();
  user.repositories.nodes.forEach(repo => {
    if (repo.isFork) return;
    repo.languages?.edges.forEach(edge => {
      const current = languageMap.get(edge.node.name) || { size: 0, color: edge.node.color };
      languageMap.set(edge.node.name, {
        size: current.size + edge.size,
        color: edge.node.color
      });
    });
  });

  const languageData = Array.from(languageMap.entries())
    .map(([name, data]) => ({ name, value: data.size, color: data.color }))
    .sort((a, b) => b.value - a.value)
    .slice(0, 5);

  const totalLoc = languageData.reduce((acc, curr) => acc + curr.value, 0);

  return (
    <main className="min-h-screen bg-[#050505] text-white relative overflow-hidden selection:bg-white/20">
      <div className="fixed inset-0 premium-grid-bg opacity-40 pointer-events-none" />
      <div className="fixed top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[400px] bg-white/5 blur-[120px] rounded-full pointer-events-none" />
      
      <div className="max-w-7xl mx-auto p-6 md:p-12 relative z-10">
        <motion.div 
          variants={container}
          initial="hidden"
          animate="show"
          className="space-y-8"
        >
          <motion.div variants={item} className="flex flex-col md:flex-row gap-10 items-start">
            <div className="relative group">
              <div className="absolute -inset-0.5 bg-gradient-to-br from-white/20 to-transparent rounded-full blur opacity-50 group-hover:opacity-100 transition duration-500" />
              <img 
                src={user.avatarUrl} 
                alt={user.login}
                className="w-32 h-32 md:w-40 md:h-40 rounded-full border-2 border-white/10 relative z-10 bg-black"
              />
            </div>
            
            <div className="flex-1 space-y-6">
              <div>
                <h1 className="text-5xl md:text-6xl font-bold tracking-tight premium-text-gradient mb-2">
                  {user.name}
                </h1>
                <a 
                  href={`https://github.com/${user.login}`}
                  target="_blank"
                  rel="noreferrer" 
                  className="text-xl text-zinc-500 font-mono hover:text-white transition-colors inline-flex items-center gap-2"
                >
                  @{user.login} <ArrowUpRight className="w-4 h-4" />
                </a>
              </div>
              
              <p className="text-lg text-zinc-400 max-w-2xl leading-relaxed font-light">
                {user.bio}
              </p>
              
              <div className="flex flex-wrap gap-6 text-sm text-zinc-500 font-medium">
                {user.location && (
                  <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/5 border border-white/5">
                    <MapPin className="w-3.5 h-3.5" /> {user.location}
                  </div>
                )}
                {user.websiteUrl && (
                  <a href={user.websiteUrl} target="_blank" rel="noreferrer" className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/5 border border-white/5 hover:bg-white/10 transition-colors">
                    <LinkIcon className="w-3.5 h-3.5" /> Website
                  </a>
                )}
                <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/5 border border-white/5">
                  <Calendar className="w-3.5 h-3.5" /> Joined {new Date(user.createdAt).getFullYear()}
                </div>
              </div>

              <div className="flex gap-8 pt-2">
                <div className="flex flex-col">
                  <span className="text-2xl font-bold text-white">{user.followers.totalCount.toLocaleString()}</span>
                  <span className="text-xs text-zinc-500 uppercase tracking-wider">Followers</span>
                </div>
                <div className="flex flex-col">
                  <span className="text-2xl font-bold text-white">{user.following.totalCount.toLocaleString()}</span>
                  <span className="text-xs text-zinc-500 uppercase tracking-wider">Following</span>
                </div>
              </div>
            </div>
          </motion.div>

          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 md:gap-6">
            <MetricCard 
              label="Total Stars" 
              value={user.repositories.nodes.reduce((acc, repo) => acc + repo.stargazers.totalCount, 0)}
              icon={<Star className="w-5 h-5 text-yellow-500" />}
              delay={0.1}
            />
            <MetricCard 
              label="Total Forks" 
              value={user.repositories.nodes.reduce((acc, repo) => acc + repo.forkCount, 0)}
              icon={<GitFork className="w-5 h-5 text-blue-500" />}
              delay={0.2}
            />
            <MetricCard 
              label="Contributions" 
              value={user.contributionsCollection.totalCommitContributions}
              icon={<Activity className="w-5 h-5 text-green-500" />}
              delay={0.3}
            />
            <MetricCard 
              label="Repositories" 
              value={user.repositories.totalCount}
              icon={<Box className="w-5 h-5 text-purple-500" />}
              delay={0.4}
            />

            <div className="md:col-span-3 grid gap-6">
              <motion.div variants={item} className="space-y-4">
                <h2 className="text-lg font-medium text-zinc-400 flex items-center gap-2">
                  <BookOpen className="w-5 h-5" /> Featured Work
                </h2>
                <div className="grid md:grid-cols-2 gap-4">
                  {pinned.user.pinnedItems.edges.map(({ node }) => (
                    <a
                      key={node.name}
                      href={node.url}
                      target="_blank"
                      rel="noreferrer"
                      className="premium-card p-6 rounded-xl flex flex-col justify-between h-full group relative overflow-hidden"
                    >
                      <div className="absolute top-0 right-0 p-4 opacity-0 group-hover:opacity-100 transition-opacity">
                        <ArrowUpRight className="w-5 h-5 text-zinc-500" />
                      </div>
                      
                      <div className="space-y-4">
                        <div className="flex items-center gap-3">
                          <div className="p-2 rounded-lg bg-white/5 border border-white/5">
                            <BookOpen className="w-5 h-5 text-zinc-300" />
                          </div>
                          <h3 className="font-bold text-lg text-zinc-100 group-hover:text-white transition-colors">
                            {node.name}
                          </h3>
                        </div>
                        <p className="text-sm text-zinc-400 leading-relaxed line-clamp-2">
                          {node.description}
                        </p>
                      </div>

                      <div className="flex items-center justify-between mt-6 pt-6 border-t border-white/5">
                        <div className="flex items-center gap-4 text-xs font-medium text-zinc-500">
                          <span className="flex items-center gap-1.5 hover:text-zinc-300 transition-colors">
                            <Star className="w-3.5 h-3.5" /> {node.stargazers.totalCount}
                          </span>
                          <span className="flex items-center gap-1.5 hover:text-zinc-300 transition-colors">
                            <GitFork className="w-3.5 h-3.5" /> {node.forkCount}
                          </span>
                        </div>
                        {node.primaryLanguage && (
                          <div className="flex items-center gap-2 text-xs text-zinc-400">
                            <span 
                              className="w-2 h-2 rounded-full shadow-[0_0_8px_rgba(0,0,0,0.5)]" 
                              style={{ backgroundColor: node.primaryLanguage.color, boxShadow: `0 0 10px ${node.primaryLanguage.color}40` }}
                            />
                            {node.primaryLanguage.name}
                          </div>
                        )}
                      </div>
                    </a>
                  ))}
                </div>
              </motion.div>
            </div>

            <div className="md:col-span-1 space-y-6">
              <motion.div variants={item} className="space-y-4">
                <h2 className="text-lg font-medium text-zinc-400 flex items-center gap-2">
                  <Layers className="w-5 h-5" /> Tech Stack
                </h2>
                <div className="premium-card p-6 rounded-xl">
                  <div className="h-[220px] w-full relative">
                    <ResponsiveContainer width="100%" height="100%">
                      <PieChart>
                        <Pie
                          data={languageData}
                          cx="50%"
                          cy="50%"
                          innerRadius={60}
                          outerRadius={80}
                          paddingAngle={4}
                          dataKey="value"
                          stroke="none"
                          cornerRadius={4}
                          onMouseEnter={(_, index) => setHoveredLang(languageData[index])}
                          onMouseLeave={() => setHoveredLang(null)}
                        >
                          {languageData.map((entry, index) => (
                            <Cell 
                              key={`cell-${index}`} 
                              fill={entry.color} 
                              className="stroke-black stroke-2 outline-none"
                            />
                          ))}
                        </Pie>
                      </PieChart>
                    </ResponsiveContainer>
                    <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none transition-opacity duration-200">
                      {hoveredLang ? (
                        <>
                          <span className="text-2xl font-bold text-white" style={{ color: hoveredLang.color }}>
                            {((hoveredLang.value / totalLoc) * 100).toFixed(1)}%
                          </span>
                          <span className="text-xs text-zinc-400 font-medium">{hoveredLang.name}</span>
                        </>
                      ) : (
                        <>
                          <span className="text-2xl font-bold text-white">{languageData.length}</span>
                          <span className="text-[10px] uppercase tracking-widest text-zinc-500">Langs</span>
                        </>
                      )}
                    </div>
                  </div>
                  
                  <div className="space-y-3 mt-6">
                    {languageData.map((lang) => (
                      <div key={lang.name} className="group flex items-center justify-between text-sm p-2 rounded-lg hover:bg-white/5 transition-colors cursor-default">
                        <div className="flex items-center gap-3">
                          <span 
                            className="w-2.5 h-2.5 rounded-full transition-transform group-hover:scale-110" 
                            style={{ backgroundColor: lang.color, boxShadow: `0 0 8px ${lang.color}40` }} 
                          />
                          <span className="text-zinc-300 font-medium">{lang.name}</span>
                        </div>
                        <span className="text-zinc-500 font-mono text-xs">
                          {((lang.value / totalLoc) * 100).toFixed(1)}%
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              </motion.div>
            </div>
          </div>
        </motion.div>
      </div>
    </main>
  );
}

function MetricCard({ label, value, icon, delay }: { label: string, value: number, icon: React.ReactNode, delay: number }) {
  return (
    <motion.div 
      variants={item}
      className="premium-card p-6 rounded-xl flex flex-col justify-between h-32 relative overflow-hidden group"
    >
      <div className="absolute top-0 right-0 p-3 opacity-10 group-hover:opacity-20 transition-opacity transform group-hover:scale-110 duration-500">
        <div className="scale-[2.5] grayscale">{icon}</div>
      </div>
      
      <div className="flex items-center gap-3 text-zinc-400 mb-2 text-sm font-medium z-10">
        <div className="p-2 rounded-md bg-white/5 border border-white/5 group-hover:border-white/10 transition-colors">
          {icon}
        </div>
        {label}
      </div>
      <div className="text-3xl font-bold text-white tracking-tight z-10 font-mono">
        {value.toLocaleString()}
      </div>
    </motion.div>
  );
}
