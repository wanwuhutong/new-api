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
package dto

// ---------------------------------------------------------------------------
// Distributor DTOs
// ---------------------------------------------------------------------------

// CreateDistributorRequest 创建分销商请求
type CreateDistributorRequest struct {
	UserId   int    `json:"user_id" binding:"required"`
	Level    int    `json:"level" binding:"required,min=1,max=3"`
	ParentId int    `json:"parent_id"`
	Status   int    `json:"status"`
}

// UpdateDistributorRequest 更新分销商请求
type UpdateDistributorRequest struct {
	Level   int `json:"level"`
	Status  int `json:"status"`
}

// DistributorInfoResponse 分销商信息响应
type DistributorInfoResponse struct {
	Id             int    `json:"id"`
	UserId         int    `json:"user_id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	Level          int    `json:"level"`
	ParentId       int    `json:"parent_id"`
	ParentUsername string `json:"parent_username,omitempty"`
	Status         int    `json:"status"`
	TotalCommission int   `json:"total_commission"` // 累计佣金
	PendingCommission int `json:"pending_commission"` // 待确认佣金
	AvailableCommission int `json:"available_commission"` // 可提现佣金
	DirectUsers    int64  `json:"direct_users"`    // 直接邀请用户数
	ChildDistributors int `json:"child_distributors"` // 下级分销商数
	CreatedAt      int64  `json:"created_at"`
}

// DistributorListResponse 分销商列表响应
type DistributorListResponse struct {
	Distributors []*DistributorInfoResponse `json:"distributors"`
	Total        int64                      `json:"total"`
	Page         int                        `json:"page"`
	PageSize     int                        `json:"page_size"`
}

// ---------------------------------------------------------------------------
// CommissionRate DTOs
// ---------------------------------------------------------------------------

// CreateCommissionRateRequest 创建佣金比例请求
type CreateCommissionRateRequest struct {
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	TypeId     int     `json:"type_id"`
	Level1Rate float64 `json:"level1_rate" binding:"gte=0,lte=1"`
	Level2Rate float64 `json:"level2_rate" binding:"gte=0,lte=1"`
	Level3Rate float64 `json:"level3_rate" binding:"gte=0,lte=1"`
	Enabled    bool    `json:"enabled"`
}

// UpdateCommissionRateRequest 更新佣金比例请求
type UpdateCommissionRateRequest struct {
	Name       string  `json:"name"`
	Level1Rate float64 `json:"level1_rate" binding:"gte=0,lte=1"`
	Level2Rate float64 `json:"level2_rate" binding:"gte=0,lte=1"`
	Level3Rate float64 `json:"level3_rate" binding:"gte=0,lte=1"`
	Enabled    bool    `json:"enabled"`
}

// CommissionRateResponse 佣金比例响应
type CommissionRateResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	TypeId      int     `json:"type_id"`
	Level1Rate  float64 `json:"level1_rate"`
	Level2Rate  float64 `json:"level2_rate"`
	Level3Rate  float64 `json:"level3_rate"`
	Enabled     bool    `json:"enabled"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// CommissionLog DTOs
// ---------------------------------------------------------------------------

// CommissionLogResponse 佣金记录响应
type CommissionLogResponse struct {
	Id            int    `json:"id"`
	DistributorId int    `json:"distributor_id"`
	Username      string `json:"username"`       // 消费用户
	UserDisplayName string `json:"user_display_name"`
	OrderId       string `json:"order_id"`
	OrderType     string `json:"order_type"`
	Amount        int    `json:"amount"`
	Commission    int    `json:"commission"`
	Level         int    `json:"level"`
	Status        int    `json:"status"`
	Remark        string `json:"remark"`
	CreatedAt     int64  `json:"created_at"`
	SettledAt     *int64 `json:"settled_at,omitempty"`
}

// CommissionLogListResponse 佣金记录列表响应
type CommissionLogListResponse struct {
	Logs     []*CommissionLogResponse `json:"logs"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
}

// ---------------------------------------------------------------------------
// Dashboard DTOs
// ---------------------------------------------------------------------------

// DistributorDashboardResponse 分销商仪表盘响应
type DistributorDashboardResponse struct {
	TotalCommission     int `json:"total_commission"`      // 累计佣金
	PendingCommission   int `json:"pending_commission"`    // 待确认佣金
	AvailableCommission int `json:"available_commission"`  // 可提现佣金
	WithdrawnCommission int `json:"withdrawn_commission"`  // 已提现佣金
	TotalUsers          int64 `json:"total_users"`         // 邀请用户总数
	Level1Users         int64 `json:"level1_users"`        // 一级用户
	Level2Users         int64 `json:"level2_users"`        // 二级用户
	Level3Users         int64 `json:"level3_users"`        // 三级用户
	TotalDistributors   int64 `json:"total_distributors"`  // 下级分销商数
	Level1Distributors int64 `json:"level1_distributors"` // 一级分销商
	Level2Distributors int64 `json:"level2_distributors"` // 二级分销商
	Level3Distributors int64 `json:"level3_distributors"` // 三级分销商
}

// ---------------------------------------------------------------------------
// Statistics DTOs
// ---------------------------------------------------------------------------

// DistributorStatisticsResponse 分销商统计数据响应
type DistributorStatisticsResponse struct {
	TotalDistributors   int64 `json:"total_distributors"`
	ActiveDistributors  int64 `json:"active_distributors"`
	TotalCommissionPaid int64 `json:"total_commission_paid"` // 已支付佣金
	TotalCommissionPending int64 `json:"total_commission_pending"` // 待结算佣金
	TotalUsersInvited   int64 `json:"total_users_invited"`
	TopDistributors     []*DistributorInfoResponse `json:"top_distributors,omitempty"`
}
