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
package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/i18n"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// Distributor APIs (Admin)
// ---------------------------------------------------------------------------

// CreateDistributor 创建分销商
func CreateDistributor(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	var req dto.CreateDistributorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	distributor, err := service.CreateDistributor(&req)
	if err != nil {
		common.ApiError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.ApiSuccess(c, distributor)
}

// UpdateDistributor 更新分销商
func UpdateDistributor(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	var req dto.UpdateDistributorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	if err := service.UpdateDistributor(id, &req); err != nil {
		common.ApiError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.ApiSuccess(c, nil)
}

// DeleteDistributor 删除分销商
func DeleteDistributor(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	// 删除分销商(软删除)
	distributor, err := model.GetDistributorById(id)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgNotFound)
		return
	}

	distributor.Status = model.DistributorStatusDisabled
	if err := model.UpdateDistributor(distributor); err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, nil)
}

// GetAllDistributors 获取所有分销商列表(管理员用)
func GetAllDistributors(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := service.GetAllDistributorsService(page, pageSize)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, result)
}

// GetDistributor 获取分销商详情
func GetDistributor(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	info, err := service.GetDistributorInfo(id)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgNotFound)
		return
	}

	common.ApiSuccess(c, info)
}

// GetDistributorChildren 获取下级分销商
func GetDistributorChildren(c *gin.Context) {
	distributorId := GetCurrentDistributorId(c)
	if distributorId == 0 {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := service.GetChildDistributorsService(distributorId, page, pageSize)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, result)
}

// GetDistributorDirectUsers 获取分销商直接邀请的用户
func GetDistributorDirectUsers(c *gin.Context) {
	distributorId := GetCurrentDistributorId(c)
	if distributorId == 0 {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := service.GetDistributorDirectUsersService(distributorId, page, pageSize)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetDistributorStatistics 获取分销统计数据(管理员用)
func GetDistributorStatistics(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	stats, err := service.GetDistributorStatistics()
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, stats)
}

// ---------------------------------------------------------------------------
// CommissionRate APIs (Admin)
// ---------------------------------------------------------------------------

// CreateCommissionRate 创建佣金比例配置
func CreateCommissionRate(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	var req dto.CreateCommissionRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	rate, err := service.CreateCommissionRateService(&req)
	if err != nil {
		common.ApiError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.ApiSuccess(c, rate)
}

// UpdateCommissionRate 更新佣金比例配置
func UpdateCommissionRate(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	var req dto.UpdateCommissionRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	if err := service.UpdateCommissionRateService(id, &req); err != nil {
		common.ApiError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.ApiSuccess(c, nil)
}

// DeleteCommissionRate 删除佣金比例配置
func DeleteCommissionRate(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	if err := service.DeleteCommissionRateService(id); err != nil {
		common.ApiError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.ApiSuccess(c, nil)
}

// GetAllCommissionRates 获取所有佣金比例配置
func GetAllCommissionRates(c *gin.Context) {
	rates, err := service.GetAllCommissionRatesService()
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, rates)
}

// GetCommissionRateByType 根据类型获取佣金比例
func GetCommissionRateByType(c *gin.Context) {
	commissionType := c.Query("type")
	if commissionType == "" {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	typeId, _ := strconv.Atoi(c.Query("type_id"))

	rate, err := model.GetCommissionRateByType(commissionType, typeId)
	if err != nil {
		common.ApiSuccess(c, nil)
		return
	}

	common.ApiSuccess(c, rate)
}

// ---------------------------------------------------------------------------
// CommissionLog APIs
// ---------------------------------------------------------------------------

// GetDistributorCommissionLogs 获取分销商佣金记录
func GetDistributorCommissionLogs(c *gin.Context) {
	distributorId := GetCurrentDistributorId(c)
	if distributorId == 0 {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := service.GetDistributorCommissionLogsService(distributorId, page, pageSize)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, result)
}

// GetAllCommissionLogs 获取所有佣金记录(管理员用)
func GetAllCommissionLogs(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := service.GetAllCommissionLogsService(page, pageSize)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgDatabaseError)
		return
	}

	common.ApiSuccess(c, result)
}

// SettleDistributorCommission 结算分销商佣金
func SettleDistributorCommission(c *gin.Context) {
	if !IsAdmin(c) {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorI18n(c, i18n.MsgInvalidParams)
		return
	}

	if err := service.SettleCommission(id); err != nil {
		common.ApiError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.ApiSuccess(c, nil)
}

// ---------------------------------------------------------------------------
// Dashboard APIs
// ---------------------------------------------------------------------------

// GetDistributorDashboard 获取分销商仪表盘
func GetDistributorDashboard(c *gin.Context) {
	userId := GetCurrentUserId(c)
	if userId == 0 {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	dashboard, err := service.GetDistributorDashboard(userId)
	if err != nil {
		common.ApiError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.ApiSuccess(c, dashboard)
}

// GetMyDistributorInfo 获取当前用户的分销商信息
func GetMyDistributorInfo(c *gin.Context) {
	userId := GetCurrentUserId(c)
	if userId == 0 {
		common.ApiErrorI18n(c, i18n.MsgAccessDenied)
		return
	}

	info, err := service.GetDistributorByUserIdForService(userId)
	if err != nil {
		common.ApiError(c, http.StatusBadRequest, "user is not a distributor")
		return
	}

	common.ApiSuccess(c, info)
}

// ---------------------------------------------------------------------------
// Helper Functions
// ---------------------------------------------------------------------------

// GetCurrentDistributorId 获取当前用户的分销商ID
func GetCurrentDistributorId(c *gin.Context) int {
	userId := GetCurrentUserId(c)
	if userId == 0 {
		return 0
	}

	distributor, err := model.GetDistributorByUserId(userId)
	if err != nil {
		return 0
	}

	return distributor.Id
}

// IsDistributor 检查当前用户是否为分销商
func IsDistributor(c *gin.Context) bool {
	return GetCurrentDistributorId(c) > 0
}
