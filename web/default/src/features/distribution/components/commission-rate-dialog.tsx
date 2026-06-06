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
import { useCommissionRateStore } from '../stores/commission-rate-store'
import { api } from '@/lib/api'
import { handleResponse } from '@/lib/api'

interface CommissionRateDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  rateId?: number
  onSuccess?: () => void
}

export function CommissionRateDialog({
  open,
  onOpenChange,
  rateId,
  onSuccess,
}: CommissionRateDialogProps) {
  const { t } = useTranslation()
  const { create, update } = useCommissionRateStore()

  const [name, setName] = useState('')
  const [type, setType] = useState('channel')
  const [typeId, setTypeId] = useState('')
  const [level1Rate, setLevel1Rate] = useState('10')
  const [level2Rate, setLevel2Rate] = useState('5')
  const [level3Rate, setLevel3Rate] = useState('2')
  const [enabled, setEnabled] = useState(true)
  const [loading, setLoading] = useState(false)

  const isEdit = !!rateId

  useEffect(() => {
    if (open && rateId) {
      loadRate(rateId)
    } else if (open) {
      resetForm()
    }
  }, [open, rateId])

  const loadRate = async (id: number) => {
    try {
      const rates = useCommissionRateStore.getState().rates
      const rate = rates.find((r) => r.id === id)
      if (rate) {
        setName(rate.name)
        setType(rate.type)
        setTypeId(rate.type_id?.toString() || '')
        setLevel1Rate((rate.level1_rate * 100).toFixed(1))
        setLevel2Rate((rate.level2_rate * 100).toFixed(1))
        setLevel3Rate((rate.level3_rate * 100).toFixed(1))
        setEnabled(rate.enabled)
      }
    } catch (error) {
      console.error('Failed to load commission rate:', error)
    }
  }

  const resetForm = () => {
    setName('')
    setType('channel')
    setTypeId('')
    setLevel1Rate('10')
    setLevel2Rate('5')
    setLevel3Rate('2')
    setEnabled(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const data = {
        name,
        type,
        type_id: typeId ? parseInt(typeId) : undefined,
        level1_rate: parseFloat(level1Rate) / 100,
        level2_rate: parseFloat(level2Rate) / 100,
        level3_rate: parseFloat(level3Rate) / 100,
        enabled,
      }

      if (isEdit) {
        await update(rateId!, data)
      } else {
        await create(data)
      }
      onSuccess?.()
    } catch (error) {
      console.error('Failed to save commission rate:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>
            {isEdit
              ? t('commissionRate.edit')
              : t('commissionRate.create')}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">{t('commissionRate.name')}</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('commissionRate.namePlaceholder')}
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="type">{t('commissionRate.type')}</Label>
            <Select value={type} onValueChange={setType}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="channel">{t('commissionRate.type.channel')}</SelectItem>
                <SelectItem value="model">{t('commissionRate.type.model')}</SelectItem>
                <SelectItem value="product">{t('commissionRate.type.product')}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="typeId">{t('commissionRate.typeId')}</Label>
            <Input
              id="typeId"
              type="number"
              value={typeId}
              onChange={(e) => setTypeId(e.target.value)}
              placeholder={t('commissionRate.typeIdPlaceholder')}
            />
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-2">
              <Label htmlFor="level1Rate">{t('commissionRate.level1')}</Label>
              <Input
                id="level1Rate"
                type="number"
                step="0.1"
                min="0"
                max="100"
                value={level1Rate}
                onChange={(e) => setLevel1Rate(e.target.value)}
                suffix="%"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="level2Rate">{t('commissionRate.level2')}</Label>
              <Input
                id="level2Rate"
                type="number"
                step="0.1"
                min="0"
                max="100"
                value={level2Rate}
                onChange={(e) => setLevel2Rate(e.target.value)}
                suffix="%"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="level3Rate">{t('commissionRate.level3')}</Label>
              <Input
                id="level3Rate"
                type="number"
                step="0.1"
                min="0"
                max="100"
                value={level3Rate}
                onChange={(e) => setLevel3Rate(e.target.value)}
                suffix="%"
              />
            </div>
          </div>

          <div className="flex items-center justify-between">
            <Label htmlFor="enabled">{t('commissionRate.enabled')}</Label>
            <Switch
              id="enabled"
              checked={enabled}
              onCheckedChange={setEnabled}
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
