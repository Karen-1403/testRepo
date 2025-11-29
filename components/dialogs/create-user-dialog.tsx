"use client"

import type React from "react"
import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

interface CreateUserDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreateUser: (username: string, password: string, roles: string[], specificPermissions?: string) => void
}

export function CreateUserDialog({ open, onOpenChange, onCreateUser }: CreateUserDialogProps) {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [roles, setRoles] = useState("")
  const [specificPermissions, setSpecificPermissions] = useState("")
  const [permissionsError, setPermissionsError] = useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    // Validate JSON if provided
    if (specificPermissions.trim()) {
      try {
        JSON.parse(specificPermissions)
        setPermissionsError("")
      } catch {
        setPermissionsError("Invalid JSON format")
        return
      }
    }

    if (username && password) {
      const rolesArray = roles
        .split(",")
        .map((r) => r.trim())
        .filter((r) => r)
      onCreateUser(username, password, rolesArray, specificPermissions)
      setUsername("")
      setPassword("")
      setRoles("")
      setSpecificPermissions("")
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Create New User</DialogTitle>
          <DialogDescription>Add a new user account with roles and permissions</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="username">Username</Label>
            <Input
              id="username"
              placeholder="john.doe"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="roles">Roles (comma-separated)</Label>
            <Input
              id="roles"
              placeholder="admin, editor, viewer"
              value={roles}
              onChange={(e) => setRoles(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="specificPermissions">Specific Permissions (JSON)</Label>
            <textarea
              id="specificPermissions"
              placeholder={'{"read": true, "write": false, "admin": false}'}
              value={specificPermissions}
              onChange={(e) => {
                setSpecificPermissions(e.target.value)
                setPermissionsError("")
              }}
              className="w-full px-3 py-2 rounded-md border border-border bg-background text-foreground text-sm font-mono resize-vertical min-h-[100px] focus:outline-none focus:ring-2 focus:ring-primary"
            />
            {permissionsError && <p className="text-red-500 text-sm">{permissionsError}</p>}
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Create User</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
