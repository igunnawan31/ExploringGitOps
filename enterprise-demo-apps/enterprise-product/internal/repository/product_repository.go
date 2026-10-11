package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/igunnawan31/enterprise-product/internal/domain"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	Save(product *domain.Product) error
	FindAll(limit, offset int, search string) ([]domain.Product, int64, error)
	FindByID(id string) (domain.Product, error)
	FindValidationData(productIds []string) ([]domain.Product, error)
	Update(id string, updates map[string]interface{}) (domain.Product, error)
	Delete(id string) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepositoryInterface {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Save(product *domain.Product) error {
	if product.ID == "" {
		product.ID = uuid.NewString()
	}

	// GORM men-generate query : INSERT INTO users (id, name, email) VALUES (?, ?, ?)
	result := r.db.Create(product)
	return result.Error
}

func (r *ProductRepository) FindAll(limit, offset int, search string) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Model(&domain.Product{}).Preload("Category")
	if search != "" {
		query.Where("name ILIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) FindByID(id string) (domain.Product, error) {
	var product domain.Product
	// GORM men-generate query : SELECT * FROM users WHERE id = ? LIMIT 1
	result := r.db.First(&product, "id = ?", id)

	if result.Error != nil { // Jika error tidak nil
		if errors.Is(result.Error, gorm.ErrRecordNotFound) { // Jika error adalah tidak ditemukan datanya
			return domain.Product{}, errors.New("product not found")
		}

		return domain.Product{}, result.Error
	}

	return product, nil
}

func (r *ProductRepository) FindValidationData(productIds []string) ([]domain.Product, error) {
	var products []domain.Product

	err := r.db.Where("id IN ?", productIds).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Update(id string, updates map[string]interface{}) (domain.Product, error) {
	var product domain.Product

	result := r.db.First(&product, "id = ?", id)
	if result.Error != nil { // Jika error tidak nil
		if errors.Is(result.Error, gorm.ErrRecordNotFound) { // Jika error adalah tidak ditemukan datanya
			return domain.Product{}, errors.New("product not found")
		}

		return domain.Product{}, result.Error
	}

	if len(updates) == 0 {
		return product, nil // Tidak ada perubahan yang dilakukan
	}

	// GORM men-generate query : UPDATE users SET name = ?, email = ? WHERE id = ?
	resultUpdate := r.db.Model(&product).Updates(updates)
	if resultUpdate.Error != nil {
		return domain.Product{}, resultUpdate.Error
	}

	return product, nil
}

func (r *ProductRepository) Delete(id string) error {
	// GORM men-generate query: DELETE FROM users WHERE id = ?
	result := r.db.Delete(&domain.Product{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}

	return result.Error
}
