import { useTranslation } from 'react-i18next'
import type { DistributorStatistics } from '../api'

interface DistributorStatsCardProps {
  statistics: DistributorStatistics
}

export function DistributorStatsCard({ statistics }: DistributorStatsCardProps) {
  const { t } = useTranslation()

  const stats = [
    {
      label: t('distribution.stats.totalDistributors'),
      value: statistics.total_distributors,
      color: 'bg-blue-500',
    },
    {
      label: t('distribution.stats.activeDistributors'),
      value: statistics.active_distributors,
      color: 'bg-green-500',
    },
    {
      label: t('distribution.stats.totalUsersInvited'),
      value: statistics.total_users_invited,
      color: 'bg-purple-500',
    },
    {
      label: t('distribution.stats.commissionPaid'),
      value: `¥${statistics.total_commission_paid.toFixed(2)}`,
      color: 'bg-yellow-500',
    },
    {
      label: t('distribution.stats.commissionPending'),
      value: `¥${statistics.total_commission_pending.toFixed(2)}`,
      color: 'bg-orange-500',
    },
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
      {stats.map((stat) => (
        <div
          key={stat.label}
          className="bg-card rounded-lg border p-4 shadow-sm"
        >
          <div className="flex items-center gap-3">
            <div className={`w-2 h-10 rounded-full ${stat.color}`} />
            <div>
              <p className="text-sm text-muted-foreground">{stat.label}</p>
              <p className="text-2xl font-bold">{stat.value}</p>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
