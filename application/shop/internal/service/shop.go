package service

import (
	"github.com/go-kratos/kratos/v2/log"
	v1 "shop/api/shop/v1"
	"shop/internal/biz"
)

// ShopService is a shop service.
type ShopService struct {
	v1.UnimplementedShopServiceServer

	uc  *biz.UserUsecase
	log *log.Helper
}

// NewShopService new a shop service.
func NewShopService(uc *biz.UserUsecase, logger log.Logger) *ShopService {
	return &ShopService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/shop")),
	}
}
