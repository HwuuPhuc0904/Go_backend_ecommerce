# API Documentation

| Method | Path | Handler |
|--------|------|--------|

## API Routes

| Method | Path | Handler |
|--------|------|--------|
| POST | /api/v1/auth/login | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).Login-fm |
| POST | /api/v1/auth/register | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).RegisterUser-fm |
| GET | /api/v1/products | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).GetAllProducts-fm |
| GET | /api/v1/products/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).GetProductByID-fm |
| GET | /api/v1/products/category/:categoryId | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*ProductController).GetProductsByCategory-fm |
| GET | /api/v1/users/admin | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).GetAllUser-fm |
| GET | /api/v1/users/admin/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).GetUserByID-fm |
| PUT | /api/v1/users/admin/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).UpdateUser-fm |
| DELETE | /api/v1/users/admin/:id | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).DeleteUser-fm |
| PUT | /api/v1/users/change-password | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).ChangePassword-fm |
| GET | /api/v1/users/profile | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).GetProfile-fm |
| PUT | /api/v1/users/profile | GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller.(*UserController).UpdateProfileUser-fm |
