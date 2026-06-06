import { useTranslation } from 'react-i18next'
import type { Distributor } from '../api'

interface DistributorColumnsProps {
  onEdit?: (id: number) => void
  onDelete?: (id: number) => void
}

export function DistributorColumns({ onEdit, onDelete }: DistributorColumnsProps) {
  const { t } = useTranslation()

  return [
    {
      accessorKey: 'id',
      header: t('distribution.columns.id'),
    },
    {
      accessorKey: 'username',
      header: t('distribution.columns.username'),
    },
    {
      accessorKey: 'display_name',
      header: t('distribution.columns.displayName'),
    },
    {
      accessorKey: 'level',
      header: t('distribution.columns.level'),
      cell: ({ row }: { row: { original: Distributor } }) => {
        const level = row.original.level
        const labels = {
          1: t('distribution.level1'),
          2: t('distribution.level2'),
          3: t('distribution.level3'),
        }
        return <span>{labels[level as 1 | 2 | 3] || level}</span>
      },
    },
    {
      accessorKey: 'parent_username',
      header: t('distribution.columns.parent'),
      cell: ({ row }: { row: { original: Distributor } }) => {
        return row.original.parent_username || '-'
      },
    },
    {
      accessorKey: 'total_commission',
      header: t('distribution.columns.totalCommission'),
      cell: ({ row }: { row: { original: Distributor } }) => {
        return `¥${row.original.total_commission.toFixed(2)}`
      },
    },
    {
      accessorKey: 'available_commission',
      header: t('distribution.columns.availableCommission'),
      cell: ({ row }: { row: { original: Distributor } }) => {
        return `¥${row.original.available_commission.toFixed(2)}`
      },
    },
    {
      accessorKey: 'direct_users',
      header: t('distribution.columns.directUsers'),
    },
    {
      accessorKey: 'child_distributors',
      header: t('distribution.columns.childDistributors'),
    },
    {
      accessorKey: 'status',
      header: t('distribution.columns.status'),
      cell: ({ row }: { row: { original: Distributor } }) => {
        const status = row.original.status
        return (
          <span
            className={`px-2 py-1 rounded text-xs ${
              status === 1
                ? 'bg-green-100 text-green-800'
                : 'bg-gray-100 text-gray-800'
            }`}
          >
            {status === 1 ? t('distribution.active') : t('distribution.inactive')}
          </span>
        )
      },
    },
    {
      accessorKey: 'actions',
      header: t('common.actions'),
      cell: ({ row }: { row: { original: Distributor } }) => {
        return (
          <div className="flex gap-2">
            <button
              onClick={() => onEdit?.(row.original.id)}
              className="text-primary hover:underline"
            >
              {t('common.edit')}
            </button>
            <button
              onClick={() => onDelete?.(row.original.id)}
              className="text-destructive hover:underline"
            >
              {t('common.delete')}
            </button>
          </div>
        )
      },
    },
  ]
}
