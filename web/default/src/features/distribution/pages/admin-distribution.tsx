import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useBasePie } from '@/components/ui'
import { useDistributionStore } from '../stores/distribution-store'
import { DistributorColumns } from '../components/distributor-columns'
import { DistributorTable } from '../components/distributor-table'
import { DistributorDialog } from '../components/distributor-dialog'
import { DistributorStatsCard } from '../components/distributor-stats-card'

export function DistributionPage() {
  const { t } = useTranslation()
  const [createOpen, setCreateOpen] = useState(false)
  const [editId, setEditId] = useState<number | null>(null)

  const {
    distributors,
    total,
    page,
    pageSize,
    loading,
    statistics,
    fetchDistributors,
    fetchStatistics,
    setPage,
    setPageSize,
  } = useDistributionStore()

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">{t('distribution.title')}</h1>
          <p className="text-muted-foreground mt-1">
            {t('distribution.subtitle')}
          </p>
        </div>
        <Button onClick={() => setCreateOpen(true)}>
          {t('distribution.createDistributor')}
        </Button>
      </div>

      {statistics && <DistributorStatsCard statistics={statistics} />}

      <DistributorTable
        data={distributors}
        loading={loading}
        pagination={{
          page,
          pageSize,
          total,
          onPageChange: setPage,
          onPageSizeChange: setPageSize,
        }}
        onEdit={(id) => setEditId(id)}
      />

      <DistributorDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        onSuccess={() => {
          setCreateOpen(false)
          fetchDistributors()
          fetchStatistics()
        }}
      />

      <DistributorDialog
        open={!!editId}
        onOpenChange={(open) => !open && setEditId(null)}
        distributorId={editId!}
        onSuccess={() => {
          setEditId(null)
          fetchDistributors()
        }}
      />
    </div>
  )
}

// Re-export Button for convenience
import { Button } from '@/components/ui/button'
