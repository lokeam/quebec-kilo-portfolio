import { Button } from "@/shared/components/ui/button"
import { Card, CardContent } from "@/shared/components/ui/card"
import { Skeleton } from "@/shared/components/ui/skeleton"
//import Image from "next/image"

export function OnlineServicesEmptyPage() {

  return (
    <div className="container px-4 py-8 mx-auto space-y-12 md:py-12">
      {/* Hero Section */}
      <section className="grid gap-8 md:grid-cols-2 md:items-center">
        <div className="relative aspect-[4/3] md:aspect-square">
          <Skeleton className="absolute inset-0 w-full h-full rounded-lg" />
        </div>
        <div className="space-y-4">
          <h1 className="text-3xl font-bold tracking-tight md:text-4xl">
            Add your first online service
          </h1>
          <p className="text-muted-foreground">
            Start tracking your spending across all your online services and platforms.
          </p>
          <Button size="lg">Add an online service</Button>
        </div>
      </section>

      {/* Discovery Section */}
      <section className="space-y-8">
        <h2 className="text-2xl font-semibold tracking-tight">More to discover</h2>
        <div className="grid gap-6 md:grid-cols-2">
          {/* First Card */}
          <Card>
            <CardContent className="p-6">
              <div className="grid gap-6 sm:grid-cols-2 sm:items-center">
                <div className="relative aspect-video w-full">
                  <Skeleton className="absolute inset-0 w-full h-full rounded-lg" />
                </div>
                <div className="flex h-full flex-col justify-between space-y-3">
                  <h3 className="text-xl font-semibold">
                    Track spending across services & platforms
                  </h3>
                  <div className="space-y-3">
                    <p className="text-sm text-muted-foreground">
                      The spend tracking page allows you to see all your spending across your media library.
                    </p>
                    <Button variant="outline" size="sm">
                      Go to Spend Tracking
                    </Button>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Second Card */}
          <Card>
            <CardContent className="p-6">
              <div className="grid gap-6 sm:grid-cols-2 sm:items-center">
                <div className="relative aspect-video w-full">
                  <Skeleton className="absolute inset-0 w-full h-full rounded-lg" />
                </div>
                <div className="flex h-full flex-col justify-between space-y-3">
                  <h3 className="text-xl font-semibold">
                    Wishlists & Deal Tracking
                  </h3>
                  <div className="space-y-3">
                    <p className="text-sm text-muted-foreground">
                      Find the best deals and discounts for your wishlisted items
                    </p>
                    <Button variant="outline" size="sm">
                      Go to Wishlist
                    </Button>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </section>
    </div>
  )
}

