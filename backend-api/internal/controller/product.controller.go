package controller

import(
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    model "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/service"
    "math"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type ProductController struct {
	productService *service.ProductService
}

func NewProductController() *ProductController {
	return &ProductController {
		productService: service.NewProductService(),
	}
}



// GetAllProducts
func(pc * ProductController) GetAllProducts(c * gin.Context) {
    // Get parameters from request
    limitStr := c.DefaultQuery("limit", "10")
    pageStr := c.DefaultQuery("page", "1")
    shortBy := c.DefaultQuery("shortBy", "id")
    orderBy := c.DefaultQuery("shortType", "asc")

    // Convert parameters to int
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1
    }





    products, total,err := pc.productService.GetAllProducts(limit, page, shortBy, orderBy)
    
    if err != nil {
        global.Logger.Error("Failed to get products", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get products",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "pagination": gin.H{
            "total":       total,
            "limit":       limit,
            "page":        page,
            "total_pages": int(math.Ceil(float64(total) / float64(limit))),
        },
    }) 
}

// GetProductsByCategoryID
func(pc * ProductController) GetProductsByCategory(c * gin.Context){
    // Get category ID from URL
    categoryIDStr := c.Param("categoryId")
    
    limitStr := c.DefaultQuery("limit", "10")
    pageStr := c.DefaultQuery("page", "1")
    shortBy := c.DefaultQuery("shortBy", "id")
    orderBy := c.DefaultQuery("shortType", "asc")

    // Convert parameters to int
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1
    }



    // Get products by category ID
    products, total, err := pc.productService.GetProductsByCategory(categoryIDStr, limit, page, shortBy, orderBy)
    if err != nil {
        global.Logger.Error("Failed to get products by category ID", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get products by category ID",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "pagination": gin.H{
            "total":       total,
            "limit":       limit,
            "page":        page,
            "total_pages": int(math.Ceil(float64(total) / float64(limit))),
        },
    })
}

func (pc *ProductController) GetProductByID(c *gin.Context) {
	// Parse id from url
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) 

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	product, err :=  pc.productService.GetProductByID(uint(id))

	if err != nil {
		global.Logger.Error("Failed to get product", zap.Error(err), zap.String("id", idStr))
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Product not found",
        })
        return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
    // Tạm thời hardcode userID cho mục đích test
    var userID uint = 1
    
    // Bind JSON request vào modelCate
    var product model.Product
    if err := ctx.ShouldBindJSON(&product); err != nil {
        global.Logger.Error("Invalid product data", zap.Error(err))
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid product data: " + err.Error(),
        })
        return
    }
    
    // Tạo sản phẩm mới
    if err := c.productService.CreateProduct(&product, userID); err != nil {
        global.Logger.Error("Failed to create product", zap.Error(err))
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create product: " + err.Error(),
        })
        return
    }
    
    // Trả về kết quả
    ctx.JSON(http.StatusCreated, gin.H{
        "message": "Product created successfully",
        "data":    product,
    })
}
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
    // Lấy user ID từ context (đã set bởi auth middleware)
    _, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "User not authenticated",
        })
        return
    }
    
    // Parse product ID từ URL
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid product ID",
        })
        return
    }
    
    // Bind JSON request vào model
    var product model.Product
    if err := ctx.ShouldBindJSON(&product); err != nil {
        global.Logger.Error("Invalid product data", zap.Error(err))
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid product data: " + err.Error(),
        })
        return
    }
    
    // Set ID từ URL
    product.ID = uint(id)
    
    // Cập nhật sản phẩm
    if err := c.productService.UpdateProduct(&product); err != nil {
        global.Logger.Error("Failed to update product", zap.Error(err), zap.String("id", idStr))
        
        // Kiểm tra loại lỗi để trả về status code phù hợp
        if err.Error() == "product not found" {
            ctx.JSON(http.StatusNotFound, gin.H{
                "error": "Product not found",
            })
        } else if err.Error() == "you don't have permission to update this product" {
            ctx.JSON(http.StatusForbidden, gin.H{
                "error": "You don't have permission to update this product",
            })
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to update product: " + err.Error(),
            })
        }
        return
    }
    
    // Trả về kết quả
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Product updated successfully",
        "data":    product,
    })
}

// DeleteProduct xóa sản phẩm
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
    // Lấy user ID từ context (đã set bởi auth middleware)
    _, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "User not authenticated",
        })
        return
    }
    
    // Parse product ID từ URL
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid product ID",
        })
        return
    }
    
    // Xóa sản phẩm
    if err := c.productService.DeleteProductByID(uint(id)); err != nil {
        global.Logger.Error("Failed to delete product", zap.Error(err), zap.String("id", idStr))
        
        // Kiểm tra loại lỗi để trả về status code phù hợp
        if err.Error() == "product not found" {
            ctx.JSON(http.StatusNotFound, gin.H{
                "error": "Product not found",
            })
        } else if err.Error() == "you don't have permission to delete this product" {
            ctx.JSON(http.StatusForbidden, gin.H{
                "error": "You don't have permission to delete this product",
            })
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to delete product: " + err.Error(),
            })
        }
        return
    }
    
    // Trả về kết quả
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Product deleted successfully",
    })
}