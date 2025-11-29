"use client"
import Link from "next/link"
import { useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { RefreshCw, Search } from "lucide-react"

const SESSIONS = [
  {
    id: "sess_001",
    user: "john.doe",
    database: "prod-db",
    address: "192.168.1.100",
    status: "active",
    startTime: "2025-01-20 09:15:00",
    endTime: "",
  },
  {
    id: "sess_002",
    user: "jane.smith",
    database: "backup-db",
    address: "192.168.1.105",
    status: "active",
    startTime: "2025-01-20 10:22:00",
    endTime: "",
  },
  {
    id: "sess_003",
    user: "bob.wilson",
    database: "mongo",
    address: "192.168.1.110",
    status: "closed",
    startTime: "2025-01-20 08:00:00",
    endTime: "2025-01-20 09:30:00",
  },
]

export default function ActiveLoginPage() {
  const [sessions, setSessions] = useState(SESSIONS)
  const [searchQuery, setSearchQuery] = useState("")

  const handleRefresh = () => {
    console.log("[v0] Refreshing sessions...")
  }

  const handleRevokeLogin = (sessionId: string) => {
    setSessions(sessions.filter((s) => s.id !== sessionId))
  }

  const filteredSessions = sessions.filter(
    (session) =>
      session.user.toLowerCase().includes(searchQuery.toLowerCase()) ||
      session.database.toLowerCase().includes(searchQuery.toLowerCase()) ||
      session.address.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  const lastSession = filteredSessions.length > 0 ? filteredSessions[filteredSessions.length - 1] : null

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground mb-2">Active login</h2>
        <p className="text-muted-foreground">Monitor active database sessions and manage user logins</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Database Sessions</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-2 items-center">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by user, database, or address..."
                className="pl-10 bg-background border-border"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
            <Button size="sm" variant="outline" onClick={handleRefresh} className="gap-2 bg-transparent">
              <RefreshCw className="w-4 h-4" />
              Refresh
            </Button>
          </div>

          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Session ID</TableHead>
                  <TableHead>User</TableHead>
                  <TableHead>Database</TableHead>
                  <TableHead>Address</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Start Time</TableHead>
                  <TableHead>End Time</TableHead>
                  <TableHead className="w-20">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredSessions.map((session) => (
                  <TableRow key={session.id}>
                    <TableCell className="font-mono text-sm">{session.id}</TableCell>
                    <TableCell>{session.user}</TableCell>
                    <TableCell>{session.database}</TableCell>
                    <TableCell className="font-mono text-sm">{session.address}</TableCell>
                    <TableCell>
                      <Badge variant={session.status === "active" ? "default" : "secondary"}>{session.status}</Badge>
                    </TableCell>
                    <TableCell className="text-sm">{session.startTime}</TableCell>
                    <TableCell className="text-sm">{session.endTime || "-"}</TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Link href={`/admin/session-logs/${session.id}`}>
                          <Button size="sm" variant="outline">
                            View Logs
                          </Button>
                        </Link>
                        {session.status === "active" && (
                          <Button size="sm" variant="destructive" onClick={() => handleRevokeLogin(session.id)}>
                            Revoke
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
