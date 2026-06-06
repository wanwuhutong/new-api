import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import {
  getDistributorDashboard,
  getMyCommissionLogs,
  getChildDistributors,
  getDistributorDirectUsers,
  type DistributorDashboard,
  type CommissionLog,
  type Distributor,
} from '../api'

export function DistributorDashboardPage() {
  const { t } = useTranslation()
  const [dashboard, setDashboard] = useState<DistributorDashboard | null>(null)
  const [commissionLogs, setCommissionLogs] = useState<CommissionLog[]>([])
  const [childDistributors, setChildDistributors] = useState<Distributor[]>([])
  const [directUsers, setDirectUsers] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'dashboard' | 'logs' | 'children' | 'users'>('dashboard')

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    setLoading(true)
    try {
      const [dashRes, logsRes, childrenRes, usersRes] = await Promise.all([
        getDistributorDashboard(),
        getMyCommissionLogs(1, 20),
        getChildDistributors(1, 20),
        getDistributorDirectUsers(1, 20),
      ])
      setDashboard(dashRes)
      setCommissionLogs(logsRes.logs)
      setChildDistributors(childrenRes.distributors)
      setDirectUsers(usersRes.users)
    } catch (error) {
      console.error('Failed to fetch distributor data:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div>{t('common.loading')}</div>
  }

  if (!dashboard) {
    return <div>{t('distribution.notDistributor')}</div>
  }

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold">{t('distribution.myDashboard')}</h1>
        <p className="text-muted-foreground mt-1">
          {t('distribution.dashboardSubtitle')}
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          label={t('distribution.totalCommission')}
          value={`¥${dashboard.total_commission.toFixed(2)}`}
          color="bg-blue-500"
        />
        <StatCard
          label={t('distribution.availableCommission')}
          value={`¥${dashboard.available_commission.toFixed(2)}`}
          color="bg-green-500"
        />
        <StatCard
          label={t('distribution.pendingCommission')}
          value={`¥${dashboard.pending_commission.toFixed(2)}`}
          color="bg-yellow-500"
        />
        <StatCard
          label={t('distribution.withdrawnCommission')}
          value={`¥${dashboard.withdrawn_commission.toFixed(2)}`}
          color="bg-purple-500"
        />
      </div>

      {/* User Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          label={t('distribution.totalUsers')}
          value={dashboard.total_users.toString()}
          color="bg-cyan-500"
        />
        <StatCard
          label={t('distribution.level1Users')}
          value={dashboard.level1_users.toString()}
          color="bg-teal-500"
        />
        <StatCard
          label={t('distribution.level2Users')}
          value={dashboard.level2_users.toString()}
          color="bg-indigo-500"
        />
        <StatCard
          label={t('distribution.level3Users')}
          value={dashboard.level3_users.toString()}
          color="bg-violet-500"
        />
      </div>

      {/* Tabs */}
      <div className="border-b">
        <nav className="flex gap-4">
          <TabButton
            active={activeTab === 'dashboard'}
            onClick={() => setActiveTab('dashboard')}
          >
            {t('distribution.tabs.dashboard')}
          </TabButton>
          <TabButton
            active={activeTab === 'logs'}
            onClick={() => setActiveTab('logs')}
          >
            {t('distribution.tabs.commissionLogs')}
          </TabButton>
          <TabButton
            active={activeTab === 'children'}
            onClick={() => setActiveTab('children')}
          >
            {t('distribution.tabs.childDistributors')}
          </TabButton>
          <TabButton
            active={activeTab === 'users'}
            onClick={() => setActiveTab('users')}
          >
            {t('distribution.tabs.directUsers')}
          </TabButton>
        </nav>
      </div>

      {/* Tab Content */}
      <div className="bg-card rounded-lg border p-4">
        {activeTab === 'dashboard' && (
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">{t('distribution.distributorLevel')}</p>
              <p className="text-lg font-semibold">
                {dashboard.total_distributors > 0 ? t('distribution.distributor') : t('distribution.notDistributor')}
              </p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">{t('distribution.childDistributors')}</p>
              <p className="text-lg font-semibold">
                {dashboard.level1_distributors + dashboard.level2_distributors + dashboard.level3_distributors}
              </p>
            </div>
          </div>
        )}

        {activeTab === 'logs' && (
          <CommissionLogTable logs={commissionLogs} />
        )}

        {activeTab === 'children' && (
          <ChildDistributorTable distributors={childDistributors} />
        )}

        {activeTab === 'users' && (
          <DirectUsersTable users={directUsers} />
        )}
      </div>
    </div>
  )
}

function StatCard({ label, value, color }: { label: string; value: string; color: string }) {
  return (
    <div className="bg-card rounded-lg border p-4 shadow-sm">
      <div className="flex items-center gap-3">
        <div className={`w-2 h-10 rounded-full ${color}`} />
        <div>
          <p className="text-sm text-muted-foreground">{label}</p>
          <p className="text-2xl font-bold">{value}</p>
        </div>
      </div>
    </div>
  )
}

function TabButton({ active, onClick, children }: { active: boolean; onClick: () => void; children: React.ReactNode }) {
  return (
    <button
      onClick={onClick}
      className={`pb-2 px-1 text-sm font-medium border-b-2 transition-colors ${
        active
          ? 'border-primary text-primary'
          : 'border-transparent text-muted-foreground hover:text-foreground'
      }`}
    >
      {children}
    </button>
  )
}

function CommissionLogTable({ logs }: { logs: CommissionLog[] }) {
  const { t } = useTranslation()

  return (
    <table className="w-full">
      <thead>
        <tr className="text-left text-sm text-muted-foreground">
          <th className="pb-2">{t('commissionLog.columns.orderId')}</th>
          <th className="pb-2">{t('commissionLog.columns.amount')}</th>
          <th className="pb-2">{t('commissionLog.columns.commission')}</th>
          <th className="pb-2">{t('commissionLog.columns.level')}</th>
          <th className="pb-2">{t('commissionLog.columns.status')}</th>
          <th className="pb-2">{t('commissionLog.columns.createdAt')}</th>
        </tr>
      </thead>
      <tbody>
        {logs.map((log) => (
          <tr key={log.id} className="border-t">
            <td className="py-2 text-sm">{log.order_id}</td>
            <td className="py-2 text-sm">¥{log.amount.toFixed(2)}</td>
            <td className="py-2 text-sm text-green-600">+¥{log.commission.toFixed(2)}</td>
            <td className="py-2 text-sm">{t(`commissionLog.level${log.level}`)}</td>
            <td className="py-2 text-sm">
              <span className={`px-2 py-1 rounded text-xs ${
                log.status === 1 ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
              }`}>
                {log.status === 1 ? t('commissionLog.settled') : t('commissionLog.pending')}
              </span>
            </td>
            <td className="py-2 text-sm">{new Date(log.created_at * 1000).toLocaleDateString()}</td>
          </tr>
        ))}
        {logs.length === 0 && (
          <tr>
            <td colSpan={6} className="py-8 text-center text-muted-foreground">
              {t('common.noData')}
            </td>
          </tr>
        )}
      </tbody>
    </table>
  )
}

function ChildDistributorTable({ distributors }: { distributors: Distributor[] }) {
  const { t } = useTranslation()

  return (
    <table className="w-full">
      <thead>
        <tr className="text-left text-sm text-muted-foreground">
          <th className="pb-2">{t('distribution.columns.username')}</th>
          <th className="pb-2">{t('distribution.columns.level')}</th>
          <th className="pb-2">{t('distribution.columns.totalCommission')}</th>
          <th className="pb-2">{t('distribution.columns.status')}</th>
        </tr>
      </thead>
      <tbody>
        {distributors.map((d) => (
          <tr key={d.id} className="border-t">
            <td className="py-2 text-sm">{d.username}</td>
            <td className="py-2 text-sm">{t(`distribution.level${d.level}`)}</td>
            <td className="py-2 text-sm">¥{d.total_commission.toFixed(2)}</td>
            <td className="py-2 text-sm">
              <span className={`px-2 py-1 rounded text-xs ${
                d.status === 1 ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
              }`}>
                {d.status === 1 ? t('distribution.active') : t('distribution.inactive')}
              </span>
            </td>
          </tr>
        ))}
        {distributors.length === 0 && (
          <tr>
            <td colSpan={4} className="py-8 text-center text-muted-foreground">
              {t('common.noData')}
            </td>
          </tr>
        )}
      </tbody>
    </table>
  )
}

function DirectUsersTable({ users }: { users: any[] }) {
  const { t } = useTranslation()

  return (
    <table className="w-full">
      <thead>
        <tr className="text-left text-sm text-muted-foreground">
          <th className="pb-2">{t('user.columns.username')}</th>
          <th className="pb-2">{t('user.columns.email')}</th>
          <th className="pb-2">{t('user.columns.createdAt')}</th>
        </tr>
      </thead>
      <tbody>
        {users.map((user) => (
          <tr key={user.id} className="border-t">
            <td className="py-2 text-sm">{user.username}</td>
            <td className="py-2 text-sm">{user.email || '-'}</td>
            <td className="py-2 text-sm">{new Date(user.created_at * 1000).toLocaleDateString()}</td>
          </tr>
        ))}
        {users.length === 0 && (
          <tr>
            <td colSpan={3} className="py-8 text-center text-muted-foreground">
              {t('common.noData')}
            </td>
          </tr>
        )}
      </tbody>
    </table>
  )
}
