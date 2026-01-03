"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import Cookies from "js-cookie"
import { Loader2, Plus, Building2 } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { API_BASE_url } from "@/lib/constants"

interface Business {
  id: string
  name: string
  // Add other fields as needed based on backend response
}

export default function BusinessSelectPage() {
  const [businesses, setBusinesses] = useState<Business[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSelecting, setIsSelecting] = useState(false)
  const router = useRouter()

  useEffect(() => {
    const fetchBusinesses = async () => {
      try {
        const response = await fetch(`${API_BASE_url}/businesses/list`, {
          credentials: "include",
        })

        if (response.ok) {
          const data = await response.json()
          setBusinesses(data.businesses || [])
        } else {
          // Handle error, maybe token expired
          if (response.status === 401) {
            router.push("/login")
          }
          console.error("Failed to fetch businesses")
        }
      } catch (error) {
        console.error(error)
        toast.error("Failed to load businesses")
      } finally {
        setIsLoading(false)
      }
    }

    fetchBusinesses()
  }, [router])

  const handleSelectBusiness = async (businessId: string) => {
    setIsSelecting(true)
    try {
      const response = await fetch(`${API_BASE_url}/businesses/select/${businessId}`, {
        method: "POST",
        credentials: "include",
      })

      if (response.ok) {
        const data = await response.json()
        toast.success("Business selected")
        router.push("/") // Redirect to dashboard
      } else {
        const data = await response.json()
        toast.error(data.message || "Failed to select business")
      }
    } catch (error) {
      console.error(error)
      toast.error("Something went wrong")
    } finally {
      setIsSelecting(false)
    }
  }

  const handleCreateBusiness = () => {
    // router.push("/business/create")
    toast.info("Create business functionality coming soon")
  }

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    )
  }

  return (
    <div className="flex min-h-screen w-full items-center justify-center p-6 md:p-10">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl">Select Business</CardTitle>
          <CardDescription>
            Choose a business to continue or create a new one
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            {businesses.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                No businesses found. Create one to get started.
              </div>
            ) : (
              businesses.map((business) => (
                <Button
                  key={business.id}
                  variant="outline"
                  className="h-16 justify-start px-4 text-left"
                  onClick={() => handleSelectBusiness(business.id)}
                  disabled={isSelecting}
                >
                  <Building2 className="mr-4 h-6 w-6" />
                  <div className="flex flex-col">
                    <span className="font-semibold">{business.name}</span>
                    <span className="text-xs text-muted-foreground">ID: {business.id}</span>
                  </div>
                </Button>
              ))
            )}
          </div>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">Or</span>
            </div>
          </div>

          <Button onClick={handleCreateBusiness} className="w-full" variant="secondary">
            <Plus className="mr-2 h-4 w-4" />
            Create New Business
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
