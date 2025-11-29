"use client"

import { Users, Database, LogIn } from "lucide-react"

export default function OverviewPage() {
  const username = typeof window !== "undefined" ? localStorage.getItem("username") || "User" : "User"
  const capitalizedUsername = username.charAt(0).toUpperCase() + username.slice(1)

  const stats = [
    { label: "Total Users", value: "3", icon: Users, color: "bg-blue-500/20 text-blue-400" },
    { label: "Databases", value: "3", icon: Database, color: "bg-purple-500/20 text-purple-400" },
    { label: "Active Logins", value: "2", icon: LogIn, color: "bg-green-500/20 text-green-400" },
  ]

  const roles = [
    { name: "admin", users: 1 },
    { name: "editor", users: 1 },
    { name: "viewer", users: 1 },
    { name: "dev", users: 1 },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-foreground mb-2">Dashboard Overview</h1>
        <p className="text-foreground/60">
          Welcome back, {capitalizedUsername}. Here's what's happening with your gateways.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {stats.map((stat, idx) => {
          const Icon = stat.icon
          return (
            <div
              key={idx}
              className="bg-card border border-border rounded-lg p-6 hover:border-primary/50 transition-colors"
            >
              <div className="flex items-center gap-4">
                <div className={`p-3 rounded-lg ${stat.color} flex-shrink-0`}>
                  <Icon className="w-6 h-6" />
                </div>
                <div className="flex flex-col">
                  <p className="text-foreground/60 text-sm font-medium">{stat.label}</p>
                  <p className="text-3xl font-bold text-foreground">{stat.value}</p>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* System Status and Roles */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
        {/* System Status */}
        <div className="bg-card border border-border rounded-lg p-6">
          <h2 className="text-lg font-bold text-foreground mb-1">System Status</h2>
          <p className="text-foreground/60 text-sm mb-4">All Systems Operational</p>
          <div className="flex items-center gap-3">
            <div className="w-3 h-3 rounded-full bg-green-500"></div>
            <span className="text-foreground text-sm font-medium">All services running smoothly</span>
          </div>
        </div>

        {/* Available Roles */}
        <div className="bg-card border border-border rounded-lg p-6">
          <h2 className="text-lg font-bold text-foreground mb-4">Available Roles</h2>
          <div className="space-y-3">
            {roles.map((role, idx) => (
              <div key={idx} className="flex items-center justify-between py-2 border-b border-border/50 last:border-0">
                <span className="text-foreground font-medium capitalize">{role.name}</span>
                <span className="text-foreground/60 text-sm">({role.users} users)</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
