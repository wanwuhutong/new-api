import { useTranslation } from 'react-i18next'
import { DataTable } from '@/components/data-table'
import type { CommissionRate } from '../api'
import { CommissionRateColumns } from './commission-rate-columns'
import { useCommissionRateStore } from '../stores/commission-rate-store'

interface CommissionRateTableProps {
  data: CommissionRate[]
  loading: boolean
  onEdit: (id: number) => void
  onDelete?: (id: number) => void
}

export function CommissionRateTable({
  data,
  loading,
  onEdit,
  onDelete,
}: CommissionRateTableProps) {
  const { t } = useTranslation()
  const { remove } = useCommissionRateStore()

  const handleDelete = async (id: number) => {
    if (window.confirm(t('commissionRate.confirmDelete'))) {
      await remove(id)
    }
  }

  const columns = CommissionRateColumns({ onEdit, onDelete: handleDelete })

  return (
    <DataTable
      columns={columns}
      data={data}
      loading={loading}
    />
  )
}
