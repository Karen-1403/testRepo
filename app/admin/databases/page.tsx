"use client"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"

const DatabaseLogos = {
  mysql: (
    <svg className="w-6 h-6" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="2" y="4" width="20" height="16" rx="2" stroke="currentColor" strokeWidth="2" />
      <path d="M6 8L12 12L18 8M6 12L12 16L18 12" stroke="currentColor" strokeWidth="1.5" />
    </svg>
  ),
  postgresql: (
    <svg className="w-6 h-6" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2" />
      <path d="M12 7C14.21 7 16 8.79 16 11" stroke="currentColor" strokeWidth="1.5" />
      <path d="M12 17C9.79 17 8 15.21 8 13" stroke="currentColor" strokeWidth="1.5" />
    </svg>
  ),
  mongodb: (
    <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm0-13c-2.76 0-5 2.24-5 5s2.24 5 5 5 5-2.24 5-5-2.24-5-5-5z" />
    </svg>
  ),
  sqlserver: (
    <svg className="w-6 h-6" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="3" y="3" width="18" height="18" rx="1" stroke="currentColor" strokeWidth="2" />
      <rect
        x="5"
        y="5"
        width="14"
        height="5"
        stroke="currentColor"
        strokeWidth="1.5"
        fill="currentColor"
        opacity="0.3"
      />
      <line x1="5" y1="12" x2="19" y2="12" stroke="currentColor" strokeWidth="1.5" />
      <rect
        x="5"
        y="14"
        width="14"
        height="5"
        stroke="currentColor"
        strokeWidth="1.5"
        fill="currentColor"
        opacity="0.3"
      />
    </svg>
  ),
}

const DATABASES = [
  {
    id: 1,
    name: "MySQL",
    description: "Open-source relational database",
    logo: DatabaseLogos.mysql,
    overview:
      "MySQL is one of the most popular open-source relational databases. It's reliable, efficient, and widely used in web applications.",
    howToConnect:
      "MySQL is a popular open-source relational database. Connect using your hostname, port 3306, and credentials.",
  },
  {
    id: 2,
    name: "PostgreSQL",
    description: "Advanced open-source database",
    logo: DatabaseLogos.postgresql,
    overview:
      "PostgreSQL is a powerful open-source object-relational database known for its advanced features and reliability.",
    howToConnect:
      "PostgreSQL is a powerful relational database. Connect using your hostname, port 5432, and database credentials.",
  },
  {
    id: 3,
    name: "MongoDB",
    description: "NoSQL document database",
    logo: DatabaseLogos.mongodb,
    overview:
      "MongoDB is a leading NoSQL document database that stores data in flexible JSON-like documents, perfect for scalable applications.",
    howToConnect:
      "MongoDB is a NoSQL document database. Connect using your connection string or hostname, port 27017, and credentials.",
  },
  {
    id: 4,
    name: "SQL Server",
    description: "Enterprise relational database",
    logo: DatabaseLogos.sqlserver,
    overview:
      "Microsoft SQL Server is an enterprise-grade relational database with powerful tools for data management and analytics.",
    howToConnect:
      "SQL Server is an enterprise relational database. Connect using your server name, port 1433, and SQL authentication credentials.",
  },
]

export default function DatabasesPage() {
  const router = useRouter()
  const [selectedDb, setSelectedDb] = useState<(typeof DATABASES)[0] | null>(null)
  const [searchQuery, setSearchQuery] = useState("")

  const filteredDatabases = DATABASES.filter(
    (db) =>
      db.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      db.description.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  const handleAddConnection = (dbId: number) => {
    setSelectedDb(null)
    router.push(`/admin/databases/${dbId}`)
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold text-foreground mb-2">Databases Available</h2>
        <p className="text-muted-foreground">Connect and manage your database instances</p>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          placeholder="Search databases..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10"
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {filteredDatabases.map((db) => (
          <Card
            key={db.id}
            onClick={() => setSelectedDb(db)}
            className="hover:border-primary transition-all cursor-pointer hover:shadow-md"
          >
            <CardHeader>
              <div className="flex items-start justify-between">
                <div className="space-y-1">
                  <CardTitle className="text-lg">{db.name}</CardTitle>
                  <CardDescription className="text-xs">{db.description}</CardDescription>
                </div>
                <div className="text-cyan-400">{db.logo}</div>
              </div>
            </CardHeader>
          </Card>
        ))}
      </div>

      {filteredDatabases.length === 0 && (
        <div className="text-center py-8">
          <p className="text-muted-foreground">No databases found matching "{searchQuery}"</p>
        </div>
      )}

      <Dialog open={!!selectedDb} onOpenChange={(open) => !open && setSelectedDb(null)}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <div className="flex items-center justify-between">
              <DialogTitle>{selectedDb?.name}</DialogTitle>
            </div>
          </DialogHeader>

          {selectedDb && (
            <div className="space-y-6">
              <div className="space-y-3">
                <h3 className="font-semibold text-foreground">Overview</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{selectedDb.overview}</p>
              </div>

              <div className="space-y-3">
                <h3 className="font-semibold text-foreground">How to Connect</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{selectedDb.howToConnect}</p>
              </div>

              <Button onClick={() => handleAddConnection(selectedDb.id)} className="w-full">
                Add database connection
              </Button>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
