"use client"
import Link from "next/link"
import { ArrowRight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ZGateLogo } from "@/components/zgate-logo"

export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Header */}
      <header className="border-b border-border">
        <div className="max-w-6xl mx-auto px-6 py-4 flex justify-between items-center">
          <div className="flex items-center gap-3">
            <ZGateLogo size={32} />
            <div className="text-2xl font-bold text-primary neon-glow-cyan">zGate</div>
          </div>
          <Link href="/login">
            <Button variant="default">Login</Button>
          </Link>
        </div>
      </header>

      {/* Hero Section */}
      <main className="flex-1 flex items-center justify-center px-6 py-24">
        <div className="text-center max-w-3xl mx-auto">
          <div className="mb-12 flex justify-center">
            <div className="relative flex items-center justify-center mb-4 animate-pulse">
              <ZGateLogo size={280} />
            </div>
          </div>

          <h1 className="text-6xl font-black tracking-tighter text-foreground mb-6 neon-glow-cyan">
            Database Access Control Platform
          </h1>
          <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
            Secure, manage, and monitor your database connections with enterprise-grade access control. Experience the
            future of database management with zGate.
          </p>
          <Link href="/login">
            <Button size="lg" className="gap-2 bg-primary text-primary-foreground hover:bg-primary/90 neon-border-cyan">
              Get Started <ArrowRight className="w-4 h-4" />
            </Button>
          </Link>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-border mt-12">
        <div className="max-w-6xl mx-auto px-6 py-8 flex justify-between items-center text-sm text-muted-foreground">
          <p>© 2025 zGate Platform</p>
          <div className="flex gap-6">
            <Link href="#" className="hover:text-primary transition-colors">
              Privacy Policy
            </Link>
            <Link href="#" className="hover:text-primary transition-colors">
              Terms of Service
            </Link>
          </div>
        </div>
      </footer>
    </div>
  )
}
