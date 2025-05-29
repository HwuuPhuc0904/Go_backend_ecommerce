# API Documentation

| Method | Path | Handler |
|--------|------|--------|

## API Routes

| Method | Path | Handler |
|--------|------|--------|
| POST | /api/v1/auth/login | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).Login-fm |
| POST | /api/v1/auth/register | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).RegisterUser-fm |
| GET | /api/v1/orders | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*OrderController).GetUserOrders-fm |
| POST | /api/v1/orders | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*OrderController).CreateOrder-fm |
| GET | /api/v1/orders/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*OrderController).GetOrderByID-fm |
| PUT | /api/v1/orders/:id/cancel | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*OrderController).CancelOrder-fm |
| POST | /api/v1/orders/:id/payment | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*OrderController).ProcessPayment-fm |
| GET | /api/v1/orders/:id/shipping | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*OrderController).GetShippingInfo-fm |
| GET | /api/v1/products | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).GetAllProducts-fm |
| GET | /api/v1/products/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).GetProductByID-fm |
| POST | /api/v1/products/:id/reviews | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).CreateProductReview-fm |
| GET | /api/v1/products/:id/reviews | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).GetProductReviews-fm |
| GET | /api/v1/products/category/:categoryId | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).GetProductsByCategory-fm |
| GET | /api/v1/users/addresses | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*AddressController).GetAddresses-fm |
| POST | /api/v1/users/addresses | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*AddressController).CreateAddress-fm |
| GET | /api/v1/users/addresses/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*AddressController).GetAddress-fm |
| PUT | /api/v1/users/addresses/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*AddressController).UpdateAddress-fm |
| DELETE | /api/v1/users/addresses/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*AddressController).DeleteAddress-fm |
| PATCH | /api/v1/users/addresses/:id/default | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*AddressController).SetDefaultAddress-fm |
| GET | /api/v1/users/admin | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).GetAllUser-fm |
| PUT | /api/v1/users/admin/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).UpdateUser-fm |
| GET | /api/v1/users/admin/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).GetUserByID-fm |
| DELETE | /api/v1/users/admin/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).DeleteUser-fm |
| PUT | /api/v1/users/change-password | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).ChangePassword-fm |
| GET | /api/v1/users/profile | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).GetProfile-fm |
| PUT | /api/v1/users/profile | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).UpdateProfileUser-fm |
