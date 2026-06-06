/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
package model

import (
	"time"

	"gorm.io/gorm"
)

// Distributor 分销商表
type Distributor struct {
	Id        int            `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId    int            `json:"user_id" gorm:"uniqueIndex;not null"`
	Level     int            `json:"level" gorm:"type:int;default:1"`      // 1=一级分销商, 2=二级分销商, 3=三级分销商
	ParentId  int            `json:"parent_id" gorm:"type:int;index"`      // 上级分销商ID
	Status    int            `json:"status" gorm:"type:int;default:1"`    // 1=启用, 0=禁用
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 设置表名
func (Distributor) TableName() string {
	return "distributors"
}

// CommissionRate 佣金比例配置表
type CommissionRate struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(128);not null"` // 配置名称，如"渠道充值佣金"
	Type      string    `json:"type" gorm:"type:varchar(32);not null"` // 配置类型: channel_topup, subscription, etc.
	TypeId    int       `json:"type_id" gorm:"type:int;default:0"`      // 关联类型ID，如特定渠道ID，0表示该类型通用
	Level1Rate float64  `json:"level1_rate" gorm:"type:decimal(5,4);default:0"` // 一级佣金比例 (0.0000-1.0000)
	Level2Rate float64  `json:"level2_rate" gorm:"type:decimal(5,4);default:0"` // 二级佣金比例
	Level3Rate float64  `json:"level3_rate" gorm:"type:decimal(5,4);default:0"` // 三级佣金比例
	Enabled   bool      `json:"enabled" gorm:"type:bool;default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 设置表名
func (CommissionRate) TableName() string {
	return "commission_rates"
}

// CommissionLog 佣金记录表
type CommissionLog struct {
	Id            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	DistributorId int       `json:"distributor_id" gorm:"index;not null"`
	UserId        int       `json:"user_id" gorm:"index;not null"`         // 消费用户
	OrderId       string    `json:"order_id" gorm:"type:varchar(64);index"` // 订单ID
	OrderType     string    `json:"order_type" gorm:"type:varchar(32)"`     // 订单类型: channel_topup, subscription, etc.
	Amount        int       `json:"amount" gorm:"type:int;default:0"`      // 消费金额(Quota)
	Commission    int       `json:"commission" gorm:"type:int;default:0"`   // 佣金金额(Quota)
	Level         int       `json:"level" gorm:"type:int;default:1"`       // 佣金层级: 1/2/3
	Status        int       `json:"status" gorm:"type:int;default:0"`      // 0=待确认, 1=已结算, 2=已提现
	Remark        string    `json:"remark" gorm:"type:varchar(255)"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	SettledAt     *int64    `json:"settled_at" gorm:"type:bigint"` // 结算时间
}

// TableName 设置表名
func (CommissionLog) TableName() string {
	return "commission_logs"
}

// DistributorStatus constants
const (
	DistributorStatusEnabled  = 1
	DistributorStatusDisabled = 0
)

// CommissionStatus constants
const (
	CommissionStatusPending   = 0 // 待确认
	CommissionStatusSettled   = 1 // 已结算(可提现)
	CommissionStatusWithdrawn = 2 // 已提现
)

// CommissionLevel constants
const (
	CommissionLevel1 = 1
	CommissionLevel2 = 2
	CommissionLevel3 = 3
)

// CommissionType constants
const (
	CommissionTypeChannelTopup  = "channel_topup"  // 渠道充值
	CommissionTypeSubscription  = "subscription"   // 订阅
	CommissionTypeRedemption    = "redemption"      // 兑换码
)

// ---------------------------------------------------------------------------
// Distributor CRUD
// ---------------------------------------------------------------------------

// GetDistributorByUserId 根据用户ID获取分销商信息
func GetDistributorByUserId(userId int) (*Distributor, error) {
	var distributor Distributor
	err := DB.First(&distributor, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

// GetDistributorById 根据ID获取分销商信息
func GetDistributorById(id int) (*Distributor, error) {
	var distributor Distributor
	err := DB.First(&distributor, id).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

// CreateDistributor 创建分销商
func CreateDistributor(distributor *Distributor) error {
	return DB.Create(distributor).Error
}

// UpdateDistributor 更新分销商
func UpdateDistributor(distributor *Distributor) error {
	return DB.Save(distributor).Error
}

// GetAllDistributors 获取所有分销商(管理员用)
func GetAllDistributors(pageInfo *PageInfo) ([]*Distributor, int64, error) {
	var distributors []*Distributor
	var total int64

	err := DB.Model(&Distributor{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (pageInfo.Page - 1) * pageInfo.PageSize
	err = DB.Offset(offset).Limit(pageInfo.PageSize).
		Order("created_at DESC").
		Find(&distributors).Error
	if err != nil {
		return nil, 0, err
	}

	return distributors, total, nil
}

// GetChildDistributors 获取下级分销商
func GetChildDistributors(parentId int, pageInfo *PageInfo) ([]*Distributor, int64, error) {
	var distributors []*Distributor
	var total int64

	err := DB.Model(&Distributor{}).Where("parent_id = ?", parentId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (pageInfo.Page - 1) * pageInfo.PageSize
	err = DB.Offset(offset).Limit(pageInfo.PageSize).
		Where("parent_id = ?", parentId).
		Order("created_at DESC").
		Find(&distributors).Error
	if err != nil {
		return nil, 0, err
	}

	return distributors, total, nil
}

// GetDistributorCommissionTotal 获取分销商累计佣金
func GetDistributorCommissionTotal(distributorId int) (total int, err error) {
	var result struct {
		Total int
	}
	err = DB.Model(&CommissionLog{}).
		Select("COALESCE(SUM(commission), 0) as total").
		Where("distributor_id = ?", distributorId).
		Scan(&result).Error
	return result.Total, err
}

// GetDistributorCommissionByStatus 获取分销商指定状态的佣金
func GetDistributorCommissionByStatus(distributorId int, status int) (total int, err error) {
	var result struct {
		Total int
	}
	err = DB.Model(&CommissionLog{}).
		Select("COALESCE(SUM(commission), 0) as total").
		Where("distributor_id = ? AND status >= ?", distributorId, status).
		Scan(&result).Error
	return result.Total, err
}

// GetDistributorDirectUsers 获取分销商直接邀请的用户
func GetDistributorDirectUsers(distributorId int, pageInfo *PageInfo) ([]*User, int64, error) {
	var distributor Distributor
	if err := DB.First(&distributor, distributorId).Error; err != nil {
		return nil, 0, err
	}

	var users []*User
	var total int64

	err := DB.Model(&User{}).Where("inviter_id = ?", distributor.UserId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (pageInfo.Page - 1) * pageInfo.PageSize
	err = DB.Offset(offset).Limit(pageInfo.PageSize).
		Where("inviter_id = ?", distributor.UserId).
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ---------------------------------------------------------------------------
// CommissionRate CRUD
// ---------------------------------------------------------------------------

// GetCommissionRateByType 根据类型获取佣金比例配置
func GetCommissionRateByType(commissionType string, typeId int) (*CommissionRate, error) {
	var rate CommissionRate
	// 先尝试查找特定ID的配置
	err := DB.First(&rate, "type = ? AND type_id = ? AND enabled = ?", commissionType, typeId, true).Error
	if err == nil {
		return &rate, nil
	}
	// 如果没找到，查找通用配置
	err = DB.First(&rate, "type = ? AND type_id = 0 AND enabled = ?", commissionType, true).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

// GetAllCommissionRates 获取所有佣金比例配置
func GetAllCommissionRates() ([]*CommissionRate, error) {
	var rates []*CommissionRate
	err := DB.Order("created_at DESC").Find(&rates).Error
	return rates, err
}

// CreateCommissionRate 创建佣金比例配置
func CreateCommissionRate(rate *CommissionRate) error {
	return DB.Create(rate).Error
}

// UpdateCommissionRate 更新佣金比例配置
func UpdateCommissionRate(rate *CommissionRate) error {
	return DB.Save(rate).Error
}

// DeleteCommissionRate 删除佣金比例配置
func DeleteCommissionRate(id int) error {
	return DB.Delete(&CommissionRate{}, id).Error
}

// ---------------------------------------------------------------------------
// CommissionLog CRUD
// ---------------------------------------------------------------------------

// CreateCommissionLog 创建佣金记录
func CreateCommissionLog(log *CommissionLog) error {
	return DB.Create(log).Error
}

// GetDistributorCommissionLogs 获取分销商佣金记录
func GetDistributorCommissionLogs(distributorId int, pageInfo *PageInfo) ([]*CommissionLog, int64, error) {
	var logs []*CommissionLog
	var total int64

	err := DB.Model(&CommissionLog{}).Where("distributor_id = ?", distributorId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (pageInfo.Page - 1) * pageInfo.PageSize
	err = DB.Offset(offset).Limit(pageInfo.PageSize).
		Where("distributor_id = ?", distributorId).
		Order("created_at DESC").
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAllCommissionLogs 获取所有佣金记录(管理员用)
func GetAllCommissionLogs(pageInfo *PageInfo) ([]*CommissionLog, int64, error) {
	var logs []*CommissionLog
	var total int64

	err := DB.Model(&CommissionLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (pageInfo.Page - 1) * pageInfo.PageSize
	err = DB.Offset(offset).Limit(pageInfo.PageSize).
		Order("created_at DESC").
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// UpdateCommissionLogStatus 更新佣金记录状态
func UpdateCommissionLogStatus(id int, status int) error {
	now := time.Now().Unix()
	return DB.Model(&CommissionLog{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"settled_at":  now,
		}).Error
}

// BatchUpdateCommissionLogStatus 批量更新佣金记录状态
func BatchUpdateCommissionLogStatus(ids []int, status int) error {
	now := time.Now().Unix()
	return DB.Model(&CommissionLog{}).Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":     status,
			"settled_at": now,
		}).Error
}

// GetUserInviterChain 获取用户的邀请链(用于计算多级佣金)
func GetUserInviterChain(userId int, maxLevel int) ([]*Distributor, error) {
	var chain []*Distributor
	currentUserId := userId

	for i := 0; i < maxLevel; i++ {
		var user User
		if err := DB.First(&user, currentUserId).Error; err != nil {
			break
		}
		if user.InviterId == 0 {
			break
		}

		distributor, err := GetDistributorByUserId(user.InviterId)
		if err != nil {
			currentUserId = user.InviterId
			continue
		}
		if distributor.Status != DistributorStatusEnabled {
			currentUserId = user.InviterId
			continue
		}

		chain = append(chain, distributor)
		currentUserId = user.InviterId
	}

	return chain, nil
}

// SettleDistributorCommission 结算分销商佣金(将待确认佣金转为可提现)
func SettleDistributorCommission(distributorId int) error {
	return DB.Model(&CommissionLog{}).
		Where("distributor_id = ? AND status = ?", distributorId, CommissionStatusPending).
		Update("status", CommissionStatusSettled).Error
}
