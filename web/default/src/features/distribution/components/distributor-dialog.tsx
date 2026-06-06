import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { useDistributionStore } from '../stores/distribution-store'
import { api } from '@/lib/api'
import { handleResponse } from '@/lib/api'

interface DistributorDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  distributorId?: number
  onSuccess?: () => void
}

export function DistributorDialog({
  open,
  onOpenChange,
  distributorId,
  onSuccess,
}: DistributorDialogProps) {
  const { t } = useTranslation()
  const { create, update } = useDistributionStore()

  const [userId, setUserId] = useState('')
  const [level, setLevel] = useState('1')
  const [parentId, setParentId] = useState('')
  const [status, setStatus] = useState(true)
  const [loading, setLoading] = useState(false)

  const isEdit = !!distributorId

  useEffect(() => {
    if (open && distributorId) {
      loadDistributor(distributorId)
    } else if (open) {
      resetForm()
    }
  }, [open, distributorId])

  const loadDistributor = async (id: number) => {
    try {
      const response = await api.get(`/distribution/admin/${id}`)
      const data = handleResponse(response)
      setUserId(data.user_id.toString())
      setLevel(data.level.toString())
      setParentId(data.parent_id?.toString() || '')
      setStatus(data.status === 1)
    } catch (error) {
      console.error('Failed to load distributor:', error)
    }
  }

  const resetForm = () => {
    setUserId('')
    setLevel('1')
    setParentId('')
    setStatus(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      if (isEdit) {
        await update(distributorId!, {
          level: parseInt(level),
          status: status ? 1 : 0,
        })
      } else {
        await create({
          user_id: parseInt(userId),
          level: parseInt(level),
          parent_id: parentId ? parseInt(parentId) : undefined,
        })
      }
      onSuccess?.()
    } catch (error) {
      console.error('Failed to save distributor:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {isEdit
              ? t('distribution.editDistributor')
              : t('distribution.createDistributor')}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          {!isEdit && (
            <div className="space-y-2">
              <Label htmlFor="userId">{t('distribution.userId')}</Label>
              <Input
                id="userId"
                type="number"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                placeholder={t('distribution.userIdPlaceholder')}
                required
              />
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="level">{t('distribution.level')}</Label>
            <Select value={level} onValueChange={setLevel}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="1">{t('distribution.level1')}</SelectItem>
                <SelectItem value="2">{t('distribution.level2')}</SelectItem>
                <SelectItem value="3">{t('distribution.level3')}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {!isEdit && (
            <div className="space-y-2">
              <Label htmlFor="parentId">{t('distribution.parentDistributor')}</Label>
              <Input
                id="parentId"
                type="number"
                value={parentId}
                onChange={(e) => setParentId(e.target.value)}
                placeholder={t('distribution.parentIdPlaceholder')}
              />
            </div>
          )}

          <div className="flex items-center justify-between">
            <Label htmlFor="status">{t('distribution.status')}</Label>
            <Switch
              id="status"
              checked={status}
              onCheckedChange={setStatus}
            />
          </div>

          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              {t('common.cancel')}
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? t('common.saving') : t('common.save')}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
