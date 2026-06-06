import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { DataTable } from '@/components/data-table'
import { getAllCommissionLogs, type CommissionLog } from '../api'

export function CommissionLogPage() {
  const { t } = useTranslation()
  const [logs, setLogs] = useState<CommissionLog[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchLogs()
  }, [page, pageSize])

  const fetchLogs = async () => {
    setLoading(true)
    try {
      const response = await getAllCommissionLogs(page, pageSize)
      setLogs(response.logs)
      setTotal(response.total)
    } catch (error) {
      console.error('Failed to fetch commission logs:', error)
    } finally {
      setLoading(false)
    }
  }

  const columns = [
    {
      accessorKey: 'id',
      header: t('commissionLog.columns.id'),
    },
    {
      accessorKey: 'username',
      header: t('commissionLog.columns.distributor'),
    },
    {
      accessorKey: 'user_display_name',
      header: t('commissionLog.columns.user'),
    },
    {
      accessorKey: 'order_id',
      header: t('commissionLog.columns.orderId'),
    },
    {
      accessorKey: 'order_type',
      header: t('commissionLog.columns.orderType'),
      cell: ({ row }: { row: { original: CommissionLog } }) => {
        return t(`commissionLog.orderType.${row.original.order_type}`)
      },
    },
    {
      accessorKey: 'amount',
      header: t('commissionLog.columns.amount'),
      cell: ({ row }: { row: { original: CommissionLog } }) => {
        return `¥${row.original.amount.toFixed(2)}`
      },
    },
    {
      accessorKey: 'commission',
      header: t('commissionLog.columns.commission'),
      cell: ({ row }: { row: { original: CommissionLog } }) => {
        return (
          <span className="text-green-600">
            +¥{row.original.commission.toFixed(2)}
          </span>
        )
      },
    },
    {
      accessorKey: 'level',
      header: t('commissionLog.columns.level'),
      cell: ({ row }: { row: { original: CommissionLog } }) => {
        return t(`commissionLog.level${row.original.level}`)
      },
    },
    {
      accessorKey: 'status',
      header: t('commissionLog.columns.status'),
      cell: ({ row }: { row: { original: CommissionLog } }) => {
        return (
          <span
            className={`px-2 py-1 rounded text-xs ${
              row.original.status === 1
                ? 'bg-green-100 text-green-800'
                : 'bg-yellow-100 text-yellow-800'
            }`}
          >
            {row.original.status === 1
              ? t('commissionLog.settled')
              : t('commissionLog.pending')}
          </span>
        )
      },
    },
    {
      accessorKey: 'created_at',
      header: t('commissionLog.columns.createdAt'),
      cell: ({ row }: { row: { original: CommissionLog } }) => {
        return new Date(row.original.created_at * 1000).toLocaleString()
      },
    },
    {
      accessorKey: 'settled_at',
      header: t('commissionLog.columns.settledAt'),
      cell: ({ row }: { row: { original: CommissionLog } }) => {
        return row.original.settled_at
          ? new Date(row.original.settled_at * 1000).toLocaleString()
          : '-'
      },
    },
  ]

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold">{t('commissionLog.title')}</h1>
        <p className="text-muted-foreground mt-1">
          {t('commissionLog.subtitle')}
        </p>
      </div>

      <DataTable
        columns={columns}
        data={logs}
        loading={loading}
        pagination={{
          page,
          pageSize,
          total,
          onPageChange: setPage,
          onPageSizeChange: setPageSize,
        }}
      />
    </div>
  )
}
