"use client"

import type React from "react"
import { useState } from "react"
import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import { Database, Users, Clock, Lock, LogOut, Menu, X, Settings, TrendingUp } from "lucide-react"
import { cn } from "@/lib/utils"
import { ZGateLogo } from "@/components/zgate-logo"

interface AdminLayoutProps {
  children: React.ReactNode
}

const NAV_ITEMS = [
  { href: "/admin/overview", label: "Overview", icon: TrendingUp },
  { href: "/admin/databases", label: "Databases Available", icon: Database },
  { href: "/admin/connected-databases", label: "Connected Databases", icon: Settings },
  { href: "/admin/users", label: "Users", icon: Users },
  { href: "/admin/sessions", label: "Active login", icon: Clock }, // Updated label to "Active login"
  { href: "/admin/queries", label: "Queries", icon: Database },
  { href: "/admin/access-control", label: "Access Control", icon: Lock },
]

export default function AdminLayout({ children }: AdminLayoutProps) {
  const pathname = usePathname()
  const router = useRouter()
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const username = typeof window !== "undefined" ? localStorage.getItem("username") || "User" : "User"

  const handleLogout = () => {
    localStorage.removeItem("auth_token")
    localStorage.removeItem("username")
    router.push("/")
  }

  const getPageName = () => {
    const segments = pathname.split("/").filter(Boolean)
    if (segments.length >= 2) {
      if (segments[1] === "databases") {
        return "Databases Available"
      }
      if (segments[1] === "sessions") {
        return "Active login"
      }
      return segments[1]
        .split("-")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ")
    }
    return "Dashboard"
  }

  return (
    <div className="flex h-screen bg-background">
      <aside
        className={cn(
          "bg-sidebar border-r border-sidebar-border transition-all duration-300 flex flex-col",
          sidebarOpen ? "w-64" : "w-20",
        )}
      >
        <div className="p-4 flex items-center justify-between">
          {sidebarOpen && (
            <div className="flex items-center gap-3">
              <ZGateLogo size={24} />
              <div className="text-lg font-bold text-primary neon-glow-cyan">zGate</div>
            </div>
          )}
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-1 hover:bg-sidebar-accent rounded transition-colors"
          >
            {sidebarOpen ? (
              <X className="w-5 h-5 text-sidebar-foreground" />
            ) : (
              <Menu className="w-5 h-5 text-sidebar-foreground" />
            )}
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
          {NAV_ITEMS.map((item) => {
            const Icon = item.icon
            const isActive = pathname === item.href
            return (
              <Link key={item.href} href={item.href}>
                <div
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-md transition-colors",
                    isActive
                      ? "bg-sidebar-primary text-sidebar-primary-foreground neon-border-cyan"
                      : "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
                  )}
                >
                  <Icon className="w-5 h-5 flex-shrink-0" />
                  {sidebarOpen && <span className="text-sm font-medium">{item.label}</span>}
                </div>
              </Link>
            )
          })}
        </nav>

        {/* Logout */}
        <div className="px-4 py-4 border-t border-sidebar-border space-y-2 mx-4">
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-md text-sidebar-foreground hover:bg-sidebar-accent hover:text-primary transition-colors text-sm">
            <Settings className="w-5 h-5 flex-shrink-0" />
            {sidebarOpen && <span>Settings</span>}
          </button>
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-3 py-2 rounded-md text-sidebar-foreground hover:bg-sidebar-accent hover:text-primary transition-colors text-sm"
          >
            <LogOut className="w-5 h-5 flex-shrink-0" />
            {sidebarOpen && <span>Logout</span>}
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="border-b border-border bg-card">
          <div className="px-6 py-3 flex justify-between items-center">
            <div className="flex items-center gap-2 text-sm">
              <span className="text-foreground font-medium">Dashboard</span>
              <span className="text-foreground/40">/</span>
              <span className="text-primary font-medium neon-glow-cyan">{getPageName()}</span>
            </div>
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-xs font-bold neon-border-cyan">
                {username.charAt(0).toUpperCase()}
              </div>
              <span className="text-sm text-foreground font-medium">{username}</span>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto p-6">{children}</main>
      </div>
    </div>
  )
}
