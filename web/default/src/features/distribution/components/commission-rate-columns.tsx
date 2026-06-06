import { useTranslation } from 'react-i18next'
import type { CommissionRate } from '../api'

interface CommissionRateColumnsProps {
  onEdit?: (id: number) => void
  onDelete?: (id: number) => void
}

export function CommissionRateColumns({ onEdit, onDelete }: CommissionRateColumnsProps) {
  const { t } = useTranslation()

  return [
    {
      accessorKey: 'name',
      header: t('commissionRate.columns.name'),
    },
    {
      accessorKey: 'type',
      header: t('commissionRate.columns.type'),
      cell: ({ row }: { row: { original: CommissionRate } }) => {
        const typeLabels: Record<string, string> = {
          channel: t('commissionRate.type.channel'),
          model: t('commissionRate.type.model'),
          product: t('commissionRate.type.product'),
        }
        return typeLabels[row.original.type] || row.original.type
      },
    },
    {
      accessorKey: 'type_id',
      header: t('commissionRate.columns.typeId'),
      cell: ({ row }: { row: { original: CommissionRate } }) => {
        return row.original.type_id || '-'
      },
    },
    {
      accessorKey: 'level1_rate',
      header: t('commissionRate.columns.level1'),
      cell: ({ row }: { row: { original: CommissionRate } }) => {
        return `${(row.original.level1_rate * 100).toFixed(1)}%`
      },
    },
    {
      accessorKey: 'level2_rate',
      header: t('commissionRate.columns.level2'),
      cell: ({ row }: { row: { original: CommissionRate } }) => {
        return `${(row.original.level2_rate * 100).toFixed(1)}%`
      },
    },
    {
      accessorKey: 'level3_rate',
      header: t('commissionRate.columns.level3'),
      cell: ({ row }: { row: { original: CommissionRate } }) => {
        return `${(row.original.level3_rate * 100).toFixed(1)}%`
      },
    },
    {
      accessorKey: 'enabled',
      header: t('commissionRate.columns.status'),
      cell: ({ row }: { row: { original: CommissionRate } }) => {
        return (
          <span
            className={`px-2 py-1 rounded text-xs ${
              row.original.enabled
                ? 'bg-green-100 text-green-800'
                : 'bg-gray-100 text-gray-800'
            }`}
          >
            {row.original.enabled
              ? t('commissionRate.enabled')
              : t('commissionRate.disabled')}
          </span>
        )
      },
    },
    {
      accessorKey: 'actions',
      header: t('common.actions'),
      cell: ({ row }: { row: { original: CommissionRate } }) => {
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
