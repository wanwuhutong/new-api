import { useTranslation } from 'react-i18next'
import { DataTable } from '@/components/data-table'
import type { Distributor } from '../api'
import { DistributorColumns } from './distributor-columns'
import { useDistributionStore } from '../stores/distribution-store'

interface DistributorTableProps {
  data: Distributor[]
  loading: boolean
  pagination: {
    page: number
    pageSize: number
    total: number
    onPageChange: (page: number) => void
    onPageSizeChange: (pageSize: number) => void
  }
  onEdit: (id: number) => void
  onDelete?: (id: number) => void
}

export function DistributorTable({
  data,
  loading,
  pagination,
  onEdit,
  onDelete,
}: DistributorTableProps) {
  const { t } = useTranslation()
  const { remove, settle } = useDistributionStore()

  const handleDelete = async (id: number) => {
    if (window.confirm(t('distribution.confirmDelete'))) {
      await remove(id)
    }
  }

  const handleSettle = async (id: number) => {
    if (window.confirm(t('distribution.confirmSettle'))) {
      await settle(id)
    }
  }

  const columns = DistributorColumns({ onEdit, onDelete: handleDelete })

  // Add settle action column
  columns.push({
    accessorKey: 'settle',
    header: t('distribution.settle'),
    cell: ({ row }: { row: { original: Distributor } }) => {
      const hasPending = row.original.pending_commission > 0
      return (
        <button
          onClick={() => handleSettle(row.original.id)}
          disabled={!hasPending}
          className={`px-2 py-1 rounded text-xs ${
            hasPending
              ? 'bg-blue-100 text-blue-800 hover:bg-blue-200'
              : 'bg-gray-100 text-gray-400 cursor-not-allowed'
          }`}
        >
          {t('distribution.settleNow')}
        </button>
      )
    },
  })

  return (
    <DataTable
      columns={columns}
      data={data}
      loading={loading}
      pagination={{
        page: pagination.page,
        pageSize: pagination.pageSize,
        total: pagination.total,
        onPageChange: pagination.onPageChange,
        onPageSizeChange: pagination.onPageSizeChange,
      }}
    />
  )
}
