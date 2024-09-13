package handlers

import "WhaleWebSite/internal/models"

type ProductContent struct {
	Info            []models.Goods
	IsAuthenticated bool
}

type ProductInfo struct {
	ProdInfo        models.Goods
	IsAuthenticated bool
}

type AboutInfo struct {
	IsAuthenticated bool
}

type ProductsInCart struct {
	CartContents []models.ProductInCart
}

type ProductAdmin struct {
	Products []models.Goods
}
