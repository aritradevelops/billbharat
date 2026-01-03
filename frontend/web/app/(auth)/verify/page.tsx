"use client"

import { useState, useEffect, Suspense } from "react"
import { useSearchParams, useRouter } from "next/navigation"
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs"
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/components/ui/input-otp"
import { API_BASE_url } from "@/lib/constants"

function VerifyContent() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const email = searchParams.get("email")

  const [isLoading, setIsLoading] = useState(false)
  const [otp, setOtp] = useState("")
  const [activeTab, setActiveTab] = useState("email")

  useEffect(() => {
    if (!email) {
      toast.error("Email is missing. Redirecting to register.")
      router.push("/register")
    }
  }, [email, router])

  const handleSendOtp = async () => {
    setIsLoading(true)
    const endpoint = activeTab === "email"
      ? "/auth/send-email-verification-request"
      : "/auth/send-phone-verification-request"

    try {
      const response = await fetch(`${API_BASE_url}${endpoint}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      })
      const data = await response.json()

      if (response.ok) {
        toast.success(`OTP sent to your ${activeTab}!`)
      } else {
        toast.error(data.message || "Failed to send OTP")
      }
    } catch (error) {
      console.error(error)
      toast.error("Something went wrong. Please try again.")
    } finally {
      setIsLoading(false)
    }
  }

  const handleVerify = async () => {
    if (otp.length !== 6) {
      toast.error("Please enter a valid 6-digit OTP")
      return
    }

    setIsLoading(true)
    const endpoint = activeTab === "email"
      ? "/auth/verify-email"
      : "/auth/verify-phone"

    try {
      const response = await fetch(`${API_BASE_url}${endpoint}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, code: otp }),
      })
      const data = await response.json()

      if (response.ok) {
        toast.success(`${activeTab === "email" ? "Email" : "Phone"} verified successfully!`)
        // Clear OTP after successful verification to allow verifying the other method or proceeding
        setOtp("")
      } else {
        toast.error(data.message || "Verification failed")
      }
    } catch (error) {
      console.error(error)
      toast.error("Something went wrong. Please try again.")
    } finally {
      setIsLoading(false)
    }
  }

  if (!email) return null

  return (
    <div className="flex min-h-screen w-full items-center justify-center p-6 md:p-10">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl">Verify Account</CardTitle>
          <CardDescription>
            Verify your email and phone number for {email}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="email" onValueChange={(val) => {
            setActiveTab(val)
            setOtp("") // Reset OTP when switching tabs
          }} className="w-full">
            <TabsList className="grid w-full grid-cols-2 mb-4">
              <TabsTrigger value="email">Email</TabsTrigger>
              <TabsTrigger value="phone">Phone</TabsTrigger>
            </TabsList>

            <div className="flex flex-col gap-4 items-center py-4">
              <div className="space-y-2 text-center">
                <p className="text-sm text-muted-foreground">
                  Enter the 6-digit code sent to your {activeTab}
                </p>
              </div>
              <InputOTP maxLength={6} value={otp} onChange={setOtp}>
                <InputOTPGroup>
                  <InputOTPSlot index={0} />
                  <InputOTPSlot index={1} />
                  <InputOTPSlot index={2} />
                </InputOTPGroup>
                <InputOTPSeparator />
                <InputOTPGroup>
                  <InputOTPSlot index={3} />
                  <InputOTPSlot index={4} />
                  <InputOTPSlot index={5} />
                </InputOTPGroup>
              </InputOTP>

              <div className="flex gap-2 w-full mt-4">
                <Button variant="outline" className="flex-1" onClick={handleSendOtp} disabled={isLoading}>
                  {isLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                  Send OTP
                </Button>
                <Button className="flex-1" onClick={handleVerify} disabled={isLoading || otp.length !== 6}>
                  {isLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                  Verify
                </Button>
              </div>
            </div>

            <TabsContent value="email">
            </TabsContent>
            <TabsContent value="phone">
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  )
}

export default function VerifyPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <VerifyContent />
    </Suspense>
  )
}
