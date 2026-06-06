import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useCommissionRateStore } from '../stores/commission-rate-store'
import { CommissionRateColumns } from '../components/commission-rate-columns'
import { CommissionRateTable } from '../components/commission-rate-table'
import { CommissionRateDialog } from '../components/commission-rate-dialog'

export function CommissionRatePage() {
  const { t } = useTranslation()
  const [createOpen, setCreateOpen] = useState(false)
  const [editId, setEditId] = useState<number | null>(null)

  const { rates, loading, fetchRates } = useCommissionRateStore()

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">{t('commissionRate.title')}</h1>
          <p className="text-muted-foreground mt-1">
            {t('commissionRate.subtitle')}
          </p>
        </div>
        <Button onClick={() => setCreateOpen(true)}>
          {t('commissionRate.create')}
        </Button>
      </div>

      <CommissionRateTable
        data={rates}
        loading={loading}
        onEdit={(id) => setEditId(id)}
      />

      <CommissionRateDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        onSuccess={() => {
          setCreateOpen(false)
          fetchRates()
        }}
      />

      <CommissionRateDialog
        open={!!editId}
        onOpenChange={(open) => !open && setEditId(null)}
        rateId={editId!}
        onSuccess={() => {
          setEditId(null)
          fetchRates()
        }}
      />
    </div>
  )
}

import { Button } from '@/components/ui/button'
