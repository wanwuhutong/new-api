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
package service

import (
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
)

// ---------------------------------------------------------------------------
// Distributor Service
// ---------------------------------------------------------------------------

// CreateDistributor 创建分销商
func CreateDistributor(req *dto.CreateDistributorRequest) (*model.Distributor, error) {
	// 检查用户是否存在
	var user model.User
	if err := model.DB.First(&user, req.UserId).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 检查是否已是分销商
	existing, err := model.GetDistributorByUserId(req.UserId)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("user is already a distributor")
	}

	distributor := &model.Distributor{
		UserId:   req.UserId,
		Level:    req.Level,
		ParentId: req.ParentId,
		Status:   req.Status,
	}
	if distributor.Status == 0 {
		distributor.Status = model.DistributorStatusEnabled
	}

	if err := model.CreateDistributor(distributor); err != nil {
		return nil, err
	}

	return distributor, nil
}

// UpdateDistributor 更新分销商
func UpdateDistributor(id int, req *dto.UpdateDistributorRequest) error {
	distributor, err := model.GetDistributorById(id)
	if err != nil {
		return fmt.Errorf("distributor not found")
	}

	if req.Level > 0 {
		distributor.Level = req.Level
	}
	if req.Status > 0 {
		distributor.Status = req.Status
	}

	return model.UpdateDistributor(distributor)
}

// GetDistributorInfo 获取分销商详细信息
func GetDistributorInfo(distributorId int) (*dto.DistributorInfoResponse, error) {
	distributor, err := model.GetDistributorById(distributorId)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := model.DB.First(&user, distributor.UserId).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 获取统计数据
	totalCommission, _ := model.GetDistributorCommissionTotal(distributorId)
	pendingCommission, _ := model.GetDistributorCommissionByStatus(distributorId, model.CommissionStatusPending)
	availableCommission, _ := model.GetDistributorCommissionByStatus(distributorId, model.CommissionStatusSettled)

	// 获取直接邀请用户数
	var directUsers int64
	model.DB.Model(&model.User{}).Where("inviter_id = ?", distributor.UserId).Count(&directUsers)

	// 获取下级分销商数
	var childDistributors int64
	model.DB.Model(&model.Distributor{}).Where("parent_id = ?", distributorId).Count(&childDistributors)

	response := &dto.DistributorInfoResponse{
		Id:                distributor.Id,
		UserId:            distributor.UserId,
		Username:          user.Username,
		DisplayName:       user.DisplayName,
		Level:             distributor.Level,
		ParentId:          distributor.ParentId,
		Status:            distributor.Status,
		TotalCommission:   totalCommission,
		PendingCommission: pendingCommission,
		AvailableCommission: availableCommission,
		DirectUsers:       directUsers,
		ChildDistributors: childDistributors,
		CreatedAt:         user.CreatedAt.Unix(),
	}

	// 获取上级分销商用户名
	if distributor.ParentId > 0 {
		if parent, err := model.GetDistributorById(distributor.ParentId); err == nil {
			if parentUser, err := model.GetUserById(parent.UserId); err == nil {
				response.ParentUsername = parentUser.Username
			}
		}
	}

	return response, nil
}

// GetDistributorByUserId 根据用户ID获取分销商信息
func GetDistributorByUserIdForService(userId int) (*dto.DistributorInfoResponse, error) {
	distributor, err := model.GetDistributorByUserId(userId)
	if err != nil {
		return nil, err
	}
	return GetDistributorInfo(distributor.Id)
}

// GetAllDistributors 获取所有分销商(管理员用)
func GetAllDistributorsService(page, pageSize int) (*dto.DistributorListResponse, error) {
	pageInfo := &common.PageInfo{Page: page, PageSize: pageSize}
	distributors, total, err := model.GetAllDistributors(pageInfo)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.DistributorInfoResponse, 0, len(distributors))
	for _, d := range distributors {
		info, err := GetDistributorInfo(d.Id)
		if err != nil {
			continue
		}
		responses = append(responses, info)
	}

	return &dto.DistributorListResponse{
		Distributors: responses,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

// GetChildDistributors 获取下级分销商
func GetChildDistributorsService(distributorId int, page, pageSize int) (*dto.DistributorListResponse, error) {
	pageInfo := &common.PageInfo{Page: page, PageSize: pageSize}
	distributors, total, err := model.GetChildDistributors(distributorId, pageInfo)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.DistributorInfoResponse, 0, len(distributors))
	for _, d := range distributors {
		info, err := GetDistributorInfo(d.Id)
		if err != nil {
			continue
		}
		responses = append(responses, info)
	}

	return &dto.DistributorListResponse{
		Distributors: responses,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

// GetDistributorDirectUsers 获取分销商直接邀请的用户
func GetDistributorDirectUsersService(distributorId int, page, pageSize int) ([]*model.User, int64, error) {
	pageInfo := &common.PageInfo{Page: page, PageSize: pageSize}
	return model.GetDistributorDirectUsers(distributorId, pageInfo)
}

// ---------------------------------------------------------------------------
// CommissionRate Service
// ---------------------------------------------------------------------------

// CreateCommissionRate 创建佣金比例配置
func CreateCommissionRateService(req *dto.CreateCommissionRateRequest) (*model.CommissionRate, error) {
	rate := &model.CommissionRate{
		Name:       req.Name,
		Type:       req.Type,
		TypeId:     req.TypeId,
		Level1Rate: req.Level1Rate,
		Level2Rate: req.Level2Rate,
		Level3Rate: req.Level3Rate,
		Enabled:    req.Enabled,
	}

	if err := model.CreateCommissionRate(rate); err != nil {
		return nil, err
	}

	return rate, nil
}

// UpdateCommissionRate 更新佣金比例配置
func UpdateCommissionRateService(id int, req *dto.UpdateCommissionRateRequest) error {
	rate, err := GetCommissionRateById(id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		rate.Name = req.Name
	}
	rate.Level1Rate = req.Level1Rate
	rate.Level2Rate = req.Level2Rate
	rate.Level3Rate = req.Level3Rate
	rate.Enabled = req.Enabled

	return model.UpdateCommissionRate(rate)
}

// GetCommissionRateById 根据ID获取佣金比例配置
func GetCommissionRateById(id int) (*model.CommissionRate, error) {
	var rate model.CommissionRate
	if err := model.DB.First(&rate, id).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

// GetAllCommissionRates 获取所有佣金比例配置
func GetAllCommissionRatesService() ([]*dto.CommissionRateResponse, error) {
	rates, err := model.GetAllCommissionRates()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.CommissionRateResponse, 0, len(rates))
	for _, r := range rates {
		responses = append(responses, &dto.CommissionRateResponse{
			Id:         r.Id,
			Name:       r.Name,
			Type:       r.Type,
			TypeId:     r.TypeId,
			Level1Rate: r.Level1Rate,
			Level2Rate: r.Level2Rate,
			Level3Rate: r.Level3Rate,
			Enabled:    r.Enabled,
			CreatedAt:  r.CreatedAt.Unix(),
			UpdatedAt:  r.UpdatedAt.Unix(),
		})
	}

	return responses, nil
}

// DeleteCommissionRate 删除佣金比例配置
func DeleteCommissionRateService(id int) error {
	return model.DeleteCommissionRate(id)
}

// ---------------------------------------------------------------------------
// CommissionLog Service
// ---------------------------------------------------------------------------

// GetDistributorCommissionLogs 获取分销商佣金记录
func GetDistributorCommissionLogsService(distributorId int, page, pageSize int) (*dto.CommissionLogListResponse, error) {
	pageInfo := &common.PageInfo{Page: page, PageSize: pageSize}
	logs, total, err := model.GetDistributorCommissionLogs(distributorId, pageInfo)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.CommissionLogResponse, 0, len(logs))
	for _, l := range logs {
		// 获取消费用户信息
		var user model.User
		model.DB.First(&user, l.UserId)

		responses = append(responses, &dto.CommissionLogResponse{
			Id:              l.Id,
			DistributorId:   l.DistributorId,
			Username:        user.Username,
			UserDisplayName: user.DisplayName,
			OrderId:         l.OrderId,
			OrderType:       l.OrderType,
			Amount:          l.Amount,
			Commission:      l.Commission,
			Level:           l.Level,
			Status:          l.Status,
			Remark:          l.Remark,
			CreatedAt:       l.CreatedAt.Unix(),
			SettledAt:       l.SettledAt,
		})
	}

	return &dto.CommissionLogListResponse{
		Logs:     responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetAllCommissionLogs 获取所有佣金记录(管理员用)
func GetAllCommissionLogsService(page, pageSize int) (*dto.CommissionLogListResponse, error) {
	pageInfo := &common.PageInfo{Page: page, PageSize: pageSize}
	logs, total, err := model.GetAllCommissionLogs(pageInfo)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.CommissionLogResponse, 0, len(logs))
	for _, l := range logs {
		// 获取消费用户信息
		var user model.User
		model.DB.First(&user, l.UserId)

		responses = append(responses, &dto.CommissionLogResponse{
			Id:              l.Id,
			DistributorId:   l.DistributorId,
			Username:        user.Username,
			UserDisplayName: user.DisplayName,
			OrderId:         l.OrderId,
			OrderType:       l.OrderType,
			Amount:          l.Amount,
			Commission:      l.Commission,
			Level:           l.Level,
			Status:          l.Status,
			Remark:          l.Remark,
			CreatedAt:       l.CreatedAt.Unix(),
			SettledAt:       l.SettledAt,
		})
	}

	return &dto.CommissionLogListResponse{
		Logs:     responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// ---------------------------------------------------------------------------
// Distribution Logic - Calculate and Record Commission
// ---------------------------------------------------------------------------

// CalculateAndRecordCommission 计算并记录佣金
// 当用户消费时调用此函数
func CalculateAndRecordCommission(userId int, amount int, orderId string, orderType string) {
	// 获取用户的邀请链
	inviterChain, err := model.GetUserInviterChain(userId, 3)
	if err != nil || len(inviterChain) == 0 {
		return
	}

	// 获取佣金比例配置
	rate, err := model.GetCommissionRateByType(orderType, 0)
	if err != nil {
		logger.LogDebug(nil, "Commission rate not found for type: %s", orderType)
		return
	}

	now := time.Now().Unix()

	// 为每一级分销商计算佣金
	for i, distributor := range inviterChain {
		if distributor.Status != model.DistributorStatusEnabled {
			continue
		}

		level := i + 1
		var commissionRate float64

		switch level {
		case 1:
			commissionRate = rate.Level1Rate
		case 2:
			commissionRate = rate.Level2Rate
		case 3:
			commissionRate = rate.Level3Rate
		}

		if commissionRate <= 0 {
			continue
		}

		commission := int(float64(amount) * commissionRate)
		if commission <= 0 {
			continue
		}

		log := &model.CommissionLog{
			DistributorId: distributor.Id,
			UserId:        userId,
			OrderId:       orderId,
			OrderType:     orderType,
			Amount:        amount,
			Commission:    commission,
			Level:         level,
			Status:        model.CommissionStatusPending,
			CreatedAt:     time.Now(),
			SettledAt:     &now,
		}

		if err := model.CreateCommissionLog(log); err != nil {
			logger.LogError(nil, "Failed to create commission log: %v", err)
		} else {
			logger.LogInfo(nil, "Commission recorded: distributor=%d, level=%d, commission=%d",
				distributor.Id, level, commission)
		}
	}
}

// SettleCommission 结算佣金(将待确认转为可提现)
func SettleCommission(distributorId int) error {
	return model.SettleDistributorCommission(distributorId)
}

// GetDistributorDashboard 获取分销商仪表盘数据
func GetDistributorDashboard(userId int) (*dto.DistributorDashboardResponse, error) {
	distributor, err := model.GetDistributorByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("user is not a distributor")
	}

	// 获取各状态佣金
	totalCommission, _ := model.GetDistributorCommissionTotal(distributor.Id)
	pendingCommission, _ := model.GetDistributorCommissionByStatus(distributor.Id, model.CommissionStatusPending)
	availableCommission, _ := model.GetDistributorCommissionByStatus(distributor.Id, model.CommissionStatusSettled)
	withdrawnCommission, _ := model.GetDistributorCommissionByStatus(distributor.Id, model.CommissionStatusWithdrawn)

	// 统计各层级用户数
	levelUsers := make(map[int]int64)
	levelDistributors := make(map[int]int64)

	// 直接邀请的用户
	var directUsers int64
	model.DB.Model(&model.User{}).Where("inviter_id = ?", distributor.UserId).Count(&directUsers)
	levelUsers[1] = directUsers

	// 遍历下级分销商统计
	var childDistributors []model.Distributor
	model.DB.Where("parent_id = ?", distributor.Id).Find(&childDistributors)

	for _, child := range childDistributors {
		// 统计该分销商直接邀请的用户
		var childUsers int64
		model.DB.Model(&model.User{}).Where("inviter_id = ?", child.UserId).Count(&childUsers)
		levelUsers[child.Level+1] += childUsers
		levelDistributors[child.Level]++
	}

	return &dto.DistributorDashboardResponse{
		TotalCommission:     totalCommission,
		PendingCommission:   pendingCommission,
		AvailableCommission: availableCommission,
		WithdrawnCommission: withdrawnCommission,
		TotalUsers:          directUsers,
		Level1Users:         levelUsers[1],
		Level2Users:         levelUsers[2],
		Level3Users:         levelUsers[3],
		TotalDistributors:   int64(len(childDistributors)),
		Level1Distributors: levelDistributors[1],
		Level2Distributors: levelDistributors[2],
		Level3Distributors: levelDistributors[3],
	}, nil
}

// GetDistributorStatistics 获取分销统计数据(管理员用)
func GetDistributorStatistics() (*dto.DistributorStatisticsResponse, error) {
	var totalDistributors int64
	model.DB.Model(&model.Distributor{}).Count(&totalDistributors)

	var activeDistributors int64
	model.DB.Model(&model.Distributor{}).Where("status = ?", model.DistributorStatusEnabled).Count(&activeDistributors)

	var totalCommissionPending int64
	model.DB.Model(&model.CommissionLog{}).Where("status = ?", model.CommissionStatusPending).
		Select("COALESCE(SUM(commission), 0)").Scan(&totalCommissionPending)

	var totalCommissionPaid int64
	model.DB.Model(&model.CommissionLog{}).Where("status >= ?", model.CommissionStatusSettled).
		Select("COALESCE(SUM(commission), 0)").Scan(&totalCommissionPaid)

	var totalUsersInvited int64
	model.DB.Model(&model.User{}).Where("inviter_id > 0").Count(&totalUsersInvited)

	return &dto.DistributorStatisticsResponse{
		TotalDistributors:     totalDistributors,
		ActiveDistributors:    activeDistributors,
		TotalCommissionPending: totalCommissionPending,
		TotalCommissionPaid:   totalCommissionPaid,
		TotalUsersInvited:     totalUsersInvited,
	}, nil
}
