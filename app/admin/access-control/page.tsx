"use client"

import { useState } from "react"
import { Plus, Edit2, Trash2, Search } from "lucide-react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { CreateAccessRoleDialog } from "@/components/dialogs/create-access-role-dialog"

const ACCESS_ROLES = [
  {
    id: 1,
    name: "admin-role",
    description: "Full administrative access",
    databases: [
      { database: "prod-db", permission: "admin" },
      { database: "backup-db", permission: "admin" },
    ],
    userCount: 5,
  },
  {
    id: 2,
    name: "read-only-role",
    description: "Read-only access to databases",
    databases: [{ database: "analytics-db", permission: "read" }],
    userCount: 12,
  },
  {
    id: 3,
    name: "developer-role",
    description: "Developer access to staging and development",
    databases: [
      { database: "Staging DB", permission: "write" },
      { database: "Development DB", permission: "write" },
    ],
    userCount: 8,
  },
]

const permissionColors: Record<string, "default" | "secondary" | "destructive"> = {
  admin: "default",
  read: "secondary",
  write: "secondary",
}

export default function AccessControlPage() {
  const [roles, setRoles] = useState(ACCESS_ROLES)
  const [showDialog, setShowDialog] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")

  const handleCreateRole = (
    name: string,
    description: string,
    databases: Array<{ database: string; permission: string }>,
  ) => {
    const newRole = {
      id: Math.max(...roles.map((r) => r.id), 0) + 1,
      name,
      description,
      databases,
      userCount: 0,
    }
    setRoles([...roles, newRole])
    setShowDialog(false)
  }

  const handleDelete = (id: number) => {
    setRoles(roles.filter((role) => role.id !== id))
  }

  const filteredRoles = roles.filter(
    (role) =>
      role.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      role.description.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-start">
        <div>
          <h2 className="text-3xl font-bold text-foreground mb-2">Access Control</h2>
          <p className="text-muted-foreground">Manage role permissions and database access</p>
        </div>
        <Button onClick={() => setShowDialog(true)} className="gap-2">
          <Plus className="w-4 h-4" />
          Add Role
        </Button>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <div className="flex items-center gap-4">
            <div>
              <CardTitle>Defined Roles</CardTitle>
              <p className="text-sm text-muted-foreground mt-1">
                {filteredRoles.length} role{filteredRoles.length !== 1 ? "s" : ""} defined
              </p>
            </div>
          </div>
          <div className="relative w-64">
            <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search roles..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-8"
            />
          </div>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Role Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Database Permissions</TableHead>
                  <TableHead>Users Count</TableHead>
                  <TableHead className="w-20">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredRoles.map((role) => (
                  <TableRow key={role.id}>
                    <TableCell className="font-medium">{role.id}</TableCell>
                    <TableCell className="font-semibold">{role.name}</TableCell>
                    <TableCell>{role.description}</TableCell>
                    <TableCell>
                      <div className="flex flex-wrap gap-2">
                        {role.databases.map((db, idx) => (
                          <Badge key={idx} variant={permissionColors[db.permission]}>
                            {db.database}: {db.permission}
                          </Badge>
                        ))}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">
                        {role.userCount} user{role.userCount !== 1 ? "s" : ""}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button size="icon" variant="ghost" className="h-8 w-8">
                          <Edit2 className="h-4 w-4" />
                        </Button>
                        <Button
                          size="icon"
                          variant="ghost"
                          className="h-8 w-8 text-destructive hover:text-destructive"
                          onClick={() => handleDelete(role.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <CreateAccessRoleDialog open={showDialog} onOpenChange={setShowDialog} onCreateRole={handleCreateRole} />
    </div>
  )
}
