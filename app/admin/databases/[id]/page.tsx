"use client"
import { useParams, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useState } from "react"

const DATABASE_TYPES = ["MySQL", "PostgreSQL", "MongoDB", "SQL Server"]
const PERMISSION_TYPES = ["Read/Write", "Read Only", "Admin", "Custom"]

export default function DatabaseDetailsPage() {
  const params = useParams()
  const router = useRouter()
  const id = params.id as string
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([])
  const [showPermissionDropdown, setShowPermissionDropdown] = useState(false)

  const [formData, setFormData] = useState({
    description: "",
    databaseName: "",
    connectionAddress: "",
    port: "",
    username: "",
    password: "",
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const databaseType = DATABASE_TYPES[Number.parseInt(id) - 1] || "Unknown"

  const handlePermissionChange = (permission: string) => {
    setSelectedPermissions((prev) =>
      prev.includes(permission) ? prev.filter((p) => p !== permission) : [...prev, permission],
    )
  }

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }))
    }
  }

  const handleConnect = async () => {
    const newErrors: Record<string, string> = {}

    if (!formData.description.trim()) newErrors.description = "Please enter in this field"
    if (!formData.databaseName.trim()) newErrors.databaseName = "Please enter in this field"
    if (!formData.connectionAddress.trim()) newErrors.connectionAddress = "Please enter in this field"
    if (!formData.port.trim()) newErrors.port = "Please enter in this field"
    if (!formData.username.trim()) newErrors.username = "Please enter in this field"
    if (!formData.password.trim()) newErrors.password = "Please enter in this field"

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors)
      return
    }

    try {
      console.log("[v0] Database connection data:", {
        databaseType,
        description: formData.description,
        databaseName: formData.databaseName,
        connectionAddress: formData.connectionAddress,
        port: formData.port,
        username: formData.username,
        permissions: selectedPermissions,
      })
      alert("Database connection configured successfully!")
      router.push("/admin/connected-databases")
    } catch (error) {
      console.error("[v0] Connection error:", error)
      alert("Failed to connect database. Please try again.")
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground mb-2">Add Database Connection</h2>
        <p className="text-muted-foreground">Configure your database connection details</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Connection Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Input
              id="description"
              placeholder="Enter database description"
              className="bg-background border-border"
              value={formData.description}
              onChange={(e) => handleInputChange("description", e.target.value)}
            />
            {errors.description && <p className="text-red-500 text-sm">{errors.description}</p>}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="type">Database Type</Label>
              <div className="bg-muted px-3 py-2 rounded-md border border-border text-foreground">{databaseType}</div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="name">Database Name</Label>
              <Input
                id="name"
                placeholder="Enter database name"
                className="bg-background border-border"
                value={formData.databaseName}
                onChange={(e) => handleInputChange("databaseName", e.target.value)}
              />
              {errors.databaseName && <p className="text-red-500 text-sm">{errors.databaseName}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="connectionAddress">Connection Address</Label>
              <Input
                id="connectionAddress"
                placeholder="e.g., 192.168.1.100"
                className="bg-background border-border"
                value={formData.connectionAddress}
                onChange={(e) => handleInputChange("connectionAddress", e.target.value)}
              />
              {errors.connectionAddress && <p className="text-red-500 text-sm">{errors.connectionAddress}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="port">Port</Label>
              <Input
                id="port"
                placeholder="e.g., 3306"
                className="bg-background border-border"
                value={formData.port}
                onChange={(e) => handleInputChange("port", e.target.value)}
              />
              {errors.port && <p className="text-red-500 text-sm">{errors.port}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                placeholder="Enter database username"
                className="bg-background border-border"
                value={formData.username}
                onChange={(e) => handleInputChange("username", e.target.value)}
              />
              {errors.username && <p className="text-red-500 text-sm">{errors.username}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="Enter database password"
                className="bg-background border-border"
                value={formData.password}
                onChange={(e) => handleInputChange("password", e.target.value)}
              />
              {errors.password && <p className="text-red-500 text-sm">{errors.password}</p>}
            </div>

            <div className="space-y-2 col-span-2">
              <Label htmlFor="defaultParams">Default Permission Type</Label>
              <button
                onClick={() => setShowPermissionDropdown(!showPermissionDropdown)}
                className="w-full bg-background border border-border rounded-md px-3 py-2 text-left text-foreground hover:border-primary transition-colors flex items-center justify-between"
              >
                <span>
                  {selectedPermissions.length === 0 ? "Select permissions..." : selectedPermissions.join(", ")}
                </span>
                <svg
                  className={`w-4 h-4 transition-transform ${showPermissionDropdown ? "rotate-180" : ""}`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
              </button>

              {showPermissionDropdown && (
                <div className="bg-background border border-t-0 border-border rounded-b-md shadow-lg max-h-40 overflow-y-auto">
                  {PERMISSION_TYPES.map((permission) => (
                    <label
                      key={permission}
                      className="flex items-center px-3 py-2 hover:bg-sidebar-accent cursor-pointer text-foreground border-b border-border last:border-b-0"
                    >
                      <input
                        type="checkbox"
                        checked={selectedPermissions.includes(permission)}
                        onChange={() => handlePermissionChange(permission)}
                        className="mr-2 w-4 h-4 rounded"
                      />
                      {permission}
                    </label>
                  ))}
                </div>
              )}
            </div>
          </div>

          <Button size="lg" className="w-full" onClick={handleConnect}>
            Connect
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
