package service

import (
	"github.com/igunnawan31/enterprise-product/internal/domain"
	"github.com/igunnawan31/enterprise-product/internal/dto"
	"github.com/igunnawan31/enterprise-product/internal/repository"
)

type ProductServiceInterface interface {
	CreateProduct(request dto.CreateProductRequest) (dto.CreateProductResponse, error)
	GetAllProduct(page, limit int, search string) ([]domain.Product, int64, error)
	GetProductById(id string) (domain.Product, error)
	GetValidationData(productIds []string) ([]dto.CreateProductResponse, error)
	UpdateProduct(id string, updates map[string]interface{}) (domain.Product, error)
	DeleteProduct(id string) error
}

type ProductService struct {
	repository repository.ProductRepositoryInterface
}

func NewProductService(repository repository.ProductRepositoryInterface) ProductServiceInterface {
	return &ProductService{repository: repository}
}

func (s *ProductService) CreateProduct(request dto.CreateProductRequest) (dto.CreateProductResponse, error) {
	product := domain.Product{
		CategoryID:  request.CategoryID,
		Name:        request.Name,
		SKU:         request.SKU,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	err := s.repository.Save(&product)
	if err != nil {
		return dto.CreateProductResponse{}, err
	}

	return dto.CreateProductResponse{}, nil
}

func (s *ProductService) GetAllProduct(page, limit int, search string) ([]domain.Product, int64, error) {
	offset := (page - 1) * limit
	return s.repository.FindAll(limit, offset, search)
}

func (s *ProductService) GetProductById(id string) (domain.Product, error) {
	return s.repository.FindByID(id)
}

func (s *ProductService) GetValidationData(productIds []string) ([]dto.CreateProductResponse, error) {
	products, err := s.repository.FindValidationData(productIds)
	if err != nil {
		return nil, err
	}

	var validProducts []dto.CreateProductResponse
	for _, product := range products {
		validProducts = append(validProducts, dto.CreateProductResponse{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		})
	}

	return validProducts, nil
}

func (s *ProductService) UpdateProduct(id string, updates map[string]interface{}) (domain.Product, error) {
	return s.repository.Update(id, updates)
}

func (s *ProductService) DeleteProduct(id string) error {
	return s.repository.Delete(id)
}
