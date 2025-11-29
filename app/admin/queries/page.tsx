"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

const CONNECTED_DATABASES = [
  { id: 1, name: "Production DB", type: "PostgreSQL" },
  { id: 2, name: "Backup Database", type: "MySQL" },
  { id: 3, name: "Analytics DB", type: "MongoDB" },
  { id: 4, name: "Staging DB", type: "PostgreSQL" },
  { id: 5, name: "Development DB", type: "MySQL" },
]

const SAMPLE_RESULTS = [
  { id: 1, name: "John Doe", email: "john@example.com", status: "active" },
  { id: 2, name: "Jane Smith", email: "jane@example.com", status: "active" },
  { id: 3, name: "Bob Wilson", email: "bob@example.com", status: "inactive" },
]

const getPlaceholder = (type: string): string => {
  switch (type.toLowerCase()) {
    case "mongodb":
      return "db.users.find({ status: 'active' })"
    case "mysql":
      return "SELECT * FROM users WHERE status = 'active'"
    case "postgresql":
      return "SELECT * FROM users WHERE status = 'active'"
    default:
      return "Enter your query here..."
  }
}

export default function QueriesPage() {
  const [selectedDb, setSelectedDb] = useState("")
  const [queryText, setQueryText] = useState("")
  const [results, setResults] = useState<typeof SAMPLE_RESULTS | null>(null)
  const [isExecuting, setIsExecuting] = useState(false)

  const handleRunQuery = () => {
    if (!selectedDb || !queryText.trim()) return
    setIsExecuting(true)
    setTimeout(() => {
      setResults(SAMPLE_RESULTS)
      setIsExecuting(false)
    }, 800)
  }

  const selectedDatabase = CONNECTED_DATABASES.find((db) => db.id.toString() === selectedDb)
  const queryLabel = selectedDatabase ? `${selectedDatabase.type} Query` : "SQL Query"
  const placeholder = selectedDatabase
    ? getPlaceholder(selectedDatabase.type)
    : "SELECT * FROM users WHERE status = 'active'"

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground mb-2">Query Runner</h2>
        <p className="text-muted-foreground">Execute queries on your connected databases</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Execute Query</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="database">Select Database</Label>
            <Select value={selectedDb} onValueChange={setSelectedDb}>
              <SelectTrigger id="database">
                <SelectValue placeholder="Choose a database" />
              </SelectTrigger>
              <SelectContent>
                {CONNECTED_DATABASES.map((db) => (
                  <SelectItem key={db.id} value={db.id.toString()}>
                    {db.name} ({db.type})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="query">{queryLabel}</Label>
            <textarea
              id="query"
              placeholder={placeholder}
              value={queryText}
              onChange={(e) => setQueryText(e.target.value)}
              className="w-full p-3 border border-border rounded-md font-mono text-xs bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary h-32"
            />
          </div>

          <Button
            onClick={handleRunQuery}
            disabled={!selectedDb || !queryText.trim() || isExecuting}
            className="w-full"
          >
            {isExecuting ? "Executing..." : "Run Query"}
          </Button>
        </CardContent>
      </Card>

      {results && (
        <Card>
          <CardHeader>
            <CardTitle>Query Results ({results.length} rows)</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>ID</TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Email</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {results.map((row) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.id}</TableCell>
                      <TableCell>{row.name}</TableCell>
                      <TableCell>{row.email}</TableCell>
                      <TableCell>{row.status}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
