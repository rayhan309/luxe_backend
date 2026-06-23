package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/luxe/backend/internal/delivery/middleware"
	jwtpkg "github.com/luxe/backend/internal/pkg/jwt"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	Auth     *AuthHandler
	Product  *ProductHandler
	Category *CategoryHandler
	Order    *OrderHandler
	Cart     *CartHandler
	Review   *ReviewHandler
	Coupon   *CouponHandler
	Banner   *BannerHandler
	User     *UserHandler
	Admin    *AdminHandler
	Upload   *UploadHandler
	JWT      *jwtpkg.Manager
}

// SetupRouter configures all routes and returns the Gin engine
func SetupRouter(deps *Dependencies, frontendURL string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(frontendURL))

	// Serve static uploaded files
	r.Static("/uploads", "./uploads")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "LUXE API"})
	})

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.Use(middleware.RateLimiter(100, 200))

	// ─── Public Auth Routes ────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Use(middleware.StrictRateLimiter())
	{
		auth.POST("/register", deps.Auth.Register)
		auth.POST("/login", deps.Auth.Login)
		auth.POST("/forgot-password", deps.Auth.ForgotPassword)
		auth.POST("/reset-password", deps.Auth.ResetPassword)
	}

	// ─── Public Product Routes ─────────────────────────────────────
	products := api.Group("/products")
	{
		products.GET("", deps.Product.List)
		products.GET("/featured", deps.Product.GetFeatured)
		products.GET("/best-sellers", deps.Product.GetBestSellers)
		products.GET("/new-arrivals", deps.Product.GetNewArrivals)
		products.GET("/:slug", deps.Product.GetBySlug)
		products.GET("/:slug/related", deps.Product.GetRelated)
	}

	// ─── Public Category Routes ────────────────────────────────────
	categories := api.Group("/categories")
	{
		categories.GET("", deps.Category.List)
		categories.GET("/tree", deps.Category.GetTree)
		categories.GET("/:slug", deps.Category.GetBySlug)
	}

	// ─── Public Banner Routes ──────────────────────────────────────
	api.GET("/banners", deps.Banner.List)

	// ─── Public Review Routes ──────────────────────────────────────
	api.GET("/reviews/:productId", deps.Review.GetProductReviews)

	// ─── Protected User Routes ─────────────────────────────────────
	authMW := middleware.Auth(deps.JWT)
	user := api.Group("/")
	user.Use(authMW)
	{
		// Profile
		user.GET("profile", deps.User.GetProfile)
		user.PUT("profile", deps.User.UpdateProfile)
		user.PUT("profile/password", deps.User.ChangePassword)

		// Addresses
		user.POST("profile/addresses", deps.User.AddAddress)
		user.PUT("profile/addresses/:addressId", deps.User.UpdateAddress)
		user.DELETE("profile/addresses/:addressId", deps.User.DeleteAddress)

		// Wishlist
		user.GET("wishlist", deps.User.GetWishlist)
		user.POST("wishlist/:productId", deps.User.ToggleWishlist)

		// Cart
		user.GET("cart", deps.Cart.GetCart)
		user.POST("cart", deps.Cart.AddItem)
		user.PUT("cart/:itemId", deps.Cart.UpdateItem)
		user.DELETE("cart/:itemId", deps.Cart.RemoveItem)
		user.DELETE("cart", deps.Cart.Clear)

		// Orders
		user.POST("orders", deps.Order.Create)
		user.GET("orders", deps.Order.GetMyOrders)
		user.GET("orders/:id", deps.Order.GetByID)

		// Reviews
		user.POST("reviews", deps.Review.Create)

		// Coupons
		user.POST("coupons/validate", deps.Coupon.Validate)
	}

	// ─── Admin Routes ──────────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(authMW, middleware.AdminOnly())
	{
		// Dashboard
		admin.GET("/dashboard", deps.Admin.Dashboard)
		admin.GET("/customers", deps.Admin.ListCustomers)

		// Products
		admin.GET("/products", deps.Product.AdminList)
		admin.POST("/products", deps.Product.AdminCreate)
		admin.PUT("/products/:id", deps.Product.AdminUpdate)
		admin.DELETE("/products/:id", deps.Product.AdminDelete)

		// Categories
		admin.GET("/categories", deps.Category.AdminList)
		admin.POST("/categories", deps.Category.AdminCreate)
		admin.PUT("/categories/:id", deps.Category.AdminUpdate)
		admin.DELETE("/categories/:id", deps.Category.AdminDelete)

		// Orders
		admin.GET("/orders", deps.Order.AdminList)
		admin.PUT("/orders/:id/status", deps.Order.AdminUpdateStatus)

		// Reviews
		admin.GET("/reviews", deps.Review.AdminList)
		admin.PUT("/reviews/:id/approve", deps.Review.AdminApprove)
		admin.DELETE("/reviews/:id", deps.Review.AdminDelete)

		// Coupons
		admin.GET("/coupons", deps.Coupon.AdminList)
		admin.POST("/coupons", deps.Coupon.AdminCreate)
		admin.PUT("/coupons/:id", deps.Coupon.AdminUpdate)
		admin.DELETE("/coupons/:id", deps.Coupon.AdminDelete)

		// Banners
		admin.GET("/banners", deps.Banner.AdminList)
		admin.POST("/banners", deps.Banner.AdminCreate)
		admin.PUT("/banners/:id", deps.Banner.AdminUpdate)
		admin.DELETE("/banners/:id", deps.Banner.AdminDelete)

		// Media upload
		admin.POST("/upload", deps.Upload.Upload)
		admin.POST("/upload/multiple", deps.Upload.UploadMultiple)
	}

	return r
}
