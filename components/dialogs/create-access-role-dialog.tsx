"use client"

import type React from "react"
import { useState } from "react"
import { Plus, X } from "lucide-react"
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
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

interface CreateAccessRoleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreateRole: (name: string, description: string, databases: Array<{ database: string; permission: string }>) => void
}

const CONNECTED_DATABASES = [
  { id: 1, name: "Production DB", type: "PostgreSQL", defaultParams: ["Read", "Write"] },
  { id: 2, name: "Backup Database", type: "MySQL", defaultParams: ["Read Only"] },
  { id: 3, name: "Analytics DB", type: "MongoDB", defaultParams: ["Admin"] },
  { id: 4, name: "Staging DB", type: "PostgreSQL", defaultParams: ["Read", "Write"] },
  { id: 5, name: "Development DB", type: "MySQL", defaultParams: ["Custom"] },
]

export function CreateAccessRoleDialog({ open, onOpenChange, onCreateRole }: CreateAccessRoleDialogProps) {
  const [roleName, setRoleName] = useState("")
  const [description, setDescription] = useState("")
  const [databases, setDatabases] = useState<Array<{ database: string; permission: string }>>([])
  const [selectedDb, setSelectedDb] = useState("")
  const [selectedPermission, setSelectedPermission] = useState("")
  const [error, setError] = useState("")

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      setRoleName("")
      setDescription("")
      setDatabases([])
      setSelectedDb("")
      setSelectedPermission("")
      setError("")
    }
    onOpenChange(newOpen)
  }

  const handleAddDatabase = () => {
    if (selectedDb && selectedPermission) {
      const dbData = CONNECTED_DATABASES.find((db) => db.id.toString() === selectedDb)
      if (dbData) {
        if (databases.some((db) => db.database === dbData.name)) {
          setError("Database permission already added")
          setTimeout(() => setError(""), 3000)
          return
        }

        setDatabases([...databases, { database: dbData.name, permission: selectedPermission }])
        setSelectedDb("")
        setSelectedPermission("")
        setError("")
      }
    }
  }

  const handleRemoveDatabase = (index: number) => {
    setDatabases(databases.filter((_, i) => i !== index))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!roleName.trim()) {
      setError("Role name is required")
      return
    }
    if (!description.trim()) {
      setError("Description is required")
      return
    }
    if (databases.length === 0) {
      setError("At least one database permission is required")
      return
    }

    onCreateRole(roleName, description, databases)
    handleOpenChange(false)
  }

  const selectedDatabase = CONNECTED_DATABASES.find((db) => db.id.toString() === selectedDb)
  const availablePermissions = selectedDatabase?.defaultParams
    ? Array.isArray(selectedDatabase.defaultParams)
      ? selectedDatabase.defaultParams
      : [selectedDatabase.defaultParams]
    : []

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Add Role</DialogTitle>
          <DialogDescription>Create a new role with database permissions</DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="roleName">Role Name</Label>
            <Input
              id="roleName"
              placeholder="e.g., admin-role, developer-role"
              value={roleName}
              onChange={(e) => setRoleName(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Input
              id="description"
              placeholder="e.g., Full administrative access"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>

          <div className="space-y-3 border-t pt-4">
            <Label>Database Permissions</Label>

            <div className="space-y-2">
              <div className="flex gap-2">
                <Select value={selectedDb} onValueChange={setSelectedDb}>
                  <SelectTrigger className="flex-1">
                    <SelectValue placeholder="Select database" />
                  </SelectTrigger>
                  <SelectContent>
                    {CONNECTED_DATABASES.map((db) => (
                      <SelectItem key={db.id} value={db.id.toString()}>
                        {db.name} ({db.type})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>

                <Select value={selectedPermission} onValueChange={setSelectedPermission} disabled={!selectedDb}>
                  <SelectTrigger className="flex-1">
                    <SelectValue placeholder="Select permission" />
                  </SelectTrigger>
                  <SelectContent>
                    {availablePermissions.length > 0 ? (
                      availablePermissions.map((perm) => (
                        <SelectItem key={perm} value={perm}>
                          {perm}
                        </SelectItem>
                      ))
                    ) : (
                      <div className="p-2 text-sm text-muted-foreground">No permissions available</div>
                    )}
                  </SelectContent>
                </Select>

                <Button type="button" size="icon" onClick={handleAddDatabase} className="gap-2">
                  <Plus className="h-4 w-4" />
                </Button>
              </div>

              {error && <p className="text-destructive text-sm">{error}</p>}
            </div>

            {databases.length > 0 && (
              <div className="space-y-2 mt-3">
                {databases.map((db, index) => (
                  <div key={index} className="flex items-center justify-between bg-muted p-2 rounded-md">
                    <span className="text-sm">
                      {db.database} - <span className="font-semibold">{db.permission}</span>
                    </span>
                    <Button
                      type="button"
                      size="icon"
                      variant="ghost"
                      className="h-6 w-6"
                      onClick={() => handleRemoveDatabase(index)}
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => handleOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit">Add Role</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
