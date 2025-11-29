"use client"

import { useState } from "react"
import { Plus, Edit2, Trash2, Search } from "lucide-react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { CreateUserDialog } from "@/components/dialogs/create-user-dialog"

const USERS = [
  {
    id: 1,
    username: "john.doe",
    roles: ["admin", "editor"],
    specificPermissions: { "read:users": true, "write:settings": false },
    createdAt: "2024-01-15",
    status: "active",
  },
  {
    id: 2,
    username: "jane.smith",
    roles: ["viewer"],
    specificPermissions: { "read:reports": true },
    createdAt: "2024-01-20",
    status: "active",
  },
  {
    id: 3,
    username: "bob.wilson",
    roles: ["editor", "analyst"],
    specificPermissions: { "read:data": true, "write:data": true, "delete:data": false },
    createdAt: "2024-02-01",
    status: "inactive",
  },
]

export default function UsersPage() {
  const [users, setUsers] = useState(USERS)
  const [showDialog, setShowDialog] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")

  const handleCreateUser = (username: string, password: string, roles: string[], specificPermissions?: string) => {
    let parsedPermissions = {}
    if (specificPermissions) {
      try {
        parsedPermissions = JSON.parse(specificPermissions)
      } catch (e) {
        parsedPermissions = {}
      }
    }

    const newUser = {
      id: Math.max(...users.map((u) => u.id), 0) + 1,
      username,
      roles,
      specificPermissions: parsedPermissions,
      createdAt: new Date().toISOString().split("T")[0],
      status: "active",
    }
    setUsers([...users, newUser])
    setShowDialog(false)
  }

  const handleDelete = (id: number) => {
    setUsers(users.filter((user) => user.id !== id))
  }

  const filteredUsers = users.filter((user) => user.username.toLowerCase().includes(searchQuery.toLowerCase()))

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-start">
        <div>
          <h2 className="text-3xl font-bold text-foreground mb-2">Users</h2>
          <p className="text-muted-foreground">Manage user accounts and permissions</p>
        </div>
        <Button onClick={() => setShowDialog(true)} className="gap-2">
          <Plus className="w-4 h-4" />
          Create User
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>User Accounts</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search users by username..."
              className="pl-10 bg-background border-border"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>

          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Username</TableHead>
                  <TableHead>Roles</TableHead>
                  <TableHead>Specific Permissions</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-20">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredUsers.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium">{user.id}</TableCell>
                    <TableCell>{user.username}</TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        {user.roles.map((role) => (
                          <span
                            key={role}
                            className="px-2 py-1 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded text-xs font-medium"
                          >
                            {role}
                          </span>
                        ))}
                      </div>
                    </TableCell>
                    <TableCell className="text-sm font-mono">
                      <code className="px-2 py-1 rounded text-xs text-muted-foreground">
                        {JSON.stringify(user.specificPermissions || {}).length > 30
                          ? JSON.stringify(user.specificPermissions || {}).substring(0, 27) + "..."
                          : JSON.stringify(user.specificPermissions || {})}
                      </code>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {new Date(user.createdAt).toLocaleDateString()}
                    </TableCell>
                    <TableCell>
                      <span
                        className={`px-2 py-1 rounded-full text-xs font-medium ${
                          user.status === "active"
                            ? "bg-green-500/10 text-green-600 dark:text-green-400"
                            : "bg-gray-500/10 text-gray-600 dark:text-gray-300"
                        }`}
                      >
                        {user.status.charAt(0).toUpperCase() + user.status.slice(1)}
                      </span>
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
                          onClick={() => handleDelete(user.id)}
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

      <CreateUserDialog open={showDialog} onOpenChange={setShowDialog} onCreateUser={handleCreateUser} />
    </div>
  )
}
