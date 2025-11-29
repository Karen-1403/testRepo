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
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

interface CreateAccessRuleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreateRule: (user: string, database: string, permission: string) => void
}

const USERS = ["john.doe", "jane.smith", "bob.wilson"]
const DATABASES = ["prod-db", "backup-db", "mongo", "analytics-db"]
const PERMISSIONS = ["admin", "read", "write"]

export function CreateAccessRuleDialog({ open, onOpenChange, onCreateRule }: CreateAccessRuleDialogProps) {
  const [user, setUser] = useState("")
  const [database, setDatabase] = useState("")
  const [permission, setPermission] = useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (user && database && permission) {
      onCreateRule(user, database, permission)
      setUser("")
      setDatabase("")
      setPermission("")
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle>Add Access Rule</DialogTitle>
          <DialogDescription>Create a new access control rule for a user and database</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="user">User</Label>
            <Select value={user} onValueChange={setUser}>
              <SelectTrigger id="user">
                <SelectValue placeholder="Select a user" />
              </SelectTrigger>
              <SelectContent>
                {USERS.map((u) => (
                  <SelectItem key={u} value={u}>
                    {u}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="database">Database</Label>
            <Select value={database} onValueChange={setDatabase}>
              <SelectTrigger id="database">
                <SelectValue placeholder="Select a database" />
              </SelectTrigger>
              <SelectContent>
                {DATABASES.map((db) => (
                  <SelectItem key={db} value={db}>
                    {db}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-2">
            <Label htmlFor="permission">Permission</Label>
            <Select value={permission} onValueChange={setPermission}>
              <SelectTrigger id="permission">
                <SelectValue placeholder="Select permission level" />
              </SelectTrigger>
              <SelectContent>
                {PERMISSIONS.map((perm) => (
                  <SelectItem key={perm} value={perm}>
                    {perm}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Add Rule</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
