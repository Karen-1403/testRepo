"use client"
import { useParams } from "next/navigation"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

const SAMPLE_LOGS = `[10:15:23] Session started for user: john.doe
[10:15:25] Connected to database: prod-db
[10:15:27] Authentication successful
[10:16:12] Query executed: SELECT * FROM users WHERE status = 'active'
[10:16:13] Query returned 2456 rows in 142ms
[10:17:45] Query executed: UPDATE users SET last_login = NOW() WHERE id = 1
[10:17:46] Update completed. Rows affected: 1
[10:18:20] Query executed: DELETE FROM sessions WHERE expires < NOW()
[10:18:21] Delete completed. Rows affected: 42
[10:19:50] Session ended
[10:19:50] Total queries executed: 3
[10:19:50] Session duration: 4m 27s`

export default function SessionLogsPage() {
  const params = useParams()
  const sessionId = params.id as string

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground mb-2">Session Logs</h2>
        <p className="text-muted-foreground">
          Session ID: <span className="font-mono">{sessionId}</span>
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Activity Log</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="bg-slate-950 text-slate-50 p-4 rounded-md font-mono text-sm leading-relaxed overflow-x-auto max-h-96 overflow-y-auto">
            {SAMPLE_LOGS.split("\n").map((line, idx) => (
              <div key={idx} className="whitespace-pre-wrap">
                {line}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
