import { ActivityFeed } from "@/components/admin/ActivityFeed"
import { RevenueChart } from "@/components/admin/RevenueChart"
import { StatCard } from "@/components/admin/StatCard"
import { dashboardStats } from "@/lib/admin-data"

export default function AdminDashboardPage() {
  return (
    <div className="space-y-6">
      <section
        aria-label="Key metrics"
        className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4"
      >
        {dashboardStats.map((stat) => (
          <StatCard key={stat.id} stat={stat} />
        ))}
      </section>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <RevenueChart />
        </div>
        <ActivityFeed />
      </div>
    </div>
  )
}
