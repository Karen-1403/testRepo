"use client"

import { useState } from "react"
import { Edit2, Trash2, Filter } from "lucide-react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuCheckboxItem,
} from "@/components/ui/dropdown-menu"

const CONNECTED_DBS = [
  {
    id: 1,
    type: "PostgreSQL",
    name: "Production DB",
    description: "Main production database",
    address: "prod-db.example.com",
    port: "5432",
    defaultParams: ["Read/Write", "Admin"],
    username: "app_user",
    createdAt: "2025-01-15 10:30:00",
  },
  {
    id: 2,
    type: "MySQL",
    name: "Backup Database",
    description: "Backup and replication database",
    address: "backup-db.example.com",
    port: "3306",
    defaultParams: ["Read Only"],
    username: "backup_user",
    createdAt: "2025-01-14 14:22:00",
  },
  {
    id: 3,
    type: "MongoDB",
    name: "Analytics DB",
    description: "Document store for analytics",
    address: "mongo.example.com",
    port: "27017",
    defaultParams: ["Admin", "Custom"],
    username: "mongo_user",
    createdAt: "2025-01-10 09:15:00",
  },
  {
    id: 4,
    type: "PostgreSQL",
    name: "Staging DB",
    description: "Staging environment database",
    address: "staging-db.example.com",
    port: "5432",
    defaultParams: ["Read/Write"],
    username: "stage_user",
    createdAt: "2025-01-09 16:45:00",
  },
  {
    id: 5,
    type: "MySQL",
    name: "Development DB",
    description: "Local development database",
    address: "dev-db.example.com",
    port: "3306",
    defaultParams: ["Read/Write", "Custom"],
    username: "dev_user",
    createdAt: "2025-01-08 11:20:00",
  },
]

const DB_TYPES = ["PostgreSQL", "MySQL", "MongoDB", "SQL Server"]

export default function ConnectedDatabasesPage() {
  const [databases, setDatabases] = useState(CONNECTED_DBS)
  const [selectedTypes, setSelectedTypes] = useState<Set<string>>(new Set(DB_TYPES))

  const handleDelete = (id: number) => {
    setDatabases(databases.filter((db) => db.id !== id))
  }

  const handleTypeToggle = (type: string) => {
    const newTypes = new Set(selectedTypes)
    if (newTypes.has(type)) {
      newTypes.delete(type)
    } else {
      newTypes.add(type)
    }
    setSelectedTypes(newTypes)
  }

  const filteredDatabases = databases.filter((db) => selectedTypes.has(db.type))

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground mb-2">Connected Databases</h2>
        <p className="text-muted-foreground">Manage your active database connections</p>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between pb-4">
          <CardTitle>Database Connections</CardTitle>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-2 bg-transparent">
                <Filter className="w-4 h-4" />
                Filter
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              {DB_TYPES.map((type) => (
                <DropdownMenuCheckboxItem
                  key={type}
                  checked={selectedTypes.has(type)}
                  onCheckedChange={() => handleTypeToggle(type)}
                >
                  {type}
                </DropdownMenuCheckboxItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Address</TableHead>
                  <TableHead>Port</TableHead>
                  <TableHead>Default Params</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead className="w-20">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredDatabases.length > 0 ? (
                  filteredDatabases.map((db) => (
                    <TableRow key={db.id}>
                      <TableCell className="font-medium">{db.id}</TableCell>
                      <TableCell className="font-medium text-primary">{db.type}</TableCell>
                      <TableCell className="font-medium">{db.name}</TableCell>
                      <TableCell className="text-sm">{db.description}</TableCell>
                      <TableCell className="text-sm">{db.address}</TableCell>
                      <TableCell className="text-sm font-medium">{db.port}</TableCell>
                      <TableCell className="text-sm">
                        {Array.isArray(db.defaultParams) ? (
                          <div className="flex flex-wrap gap-1">
                            {db.defaultParams.map((param, idx) => (
                              <span
                                key={idx}
                                className="inline-block bg-secondary text-secondary-foreground text-xs px-2 py-1 rounded"
                              >
                                {param}
                              </span>
                            ))}
                          </div>
                        ) : (
                          db.defaultParams
                        )}
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">{db.createdAt}</TableCell>
                      <TableCell>
                        <div className="flex gap-2">
                          <Button size="icon" variant="ghost" className="h-8 w-8">
                            <Edit2 className="h-4 w-4" />
                          </Button>
                          <Button
                            size="icon"
                            variant="ghost"
                            className="h-8 w-8 text-destructive hover:text-destructive"
                            onClick={() => handleDelete(db.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={9} className="text-center py-6 text-muted-foreground">
                      No databases match the selected filter
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
