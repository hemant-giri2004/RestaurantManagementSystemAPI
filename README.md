```markdown
# 🍽️ Restaurant Management System API

A RESTful API built with **Golang** and **PostgreSQL**, designed for managing restaurants, dishes, and users with **JWT-based authentication**, **role-based access control**, and **secure session management**.

> Developed by [Hemant Giri](https://github.com/hemant-giri2004)

---

## 🚀 Features

- ✅ **User Registration & Login**
  - Register with roles (Admin / Sub-Admin / User)
  - Store user address and contact info
  - Secure password hashing

- 🔐 **Authentication & Authorization**
  - JWT access token & refresh token system
  - Session-based refresh token tracking
  - Role-based access control middleware
  - Logout mechanism (token/session invalidation)

- 🧑‍💼 **Roles & Permissions**
  - Admin: Full access
  - Sub-Admin: Limited access (can manage restaurants/dishes)
  - User: View access (can view restaurants and dishes)

- 🏪 **Restaurant Management**
  - Add, update, delete restaurants (Admin/Sub-Admin)
  - View all restaurants (All roles)

- 🍔 **Dish Management**
  - Add, update, delete dishes (Admin/Sub-Admin)
  - View dishes by restaurant (All roles)

- 📦 **Modular Code Structure**
  - Handlers, DB helpers, models, middleware separated for clarity
  - Centralized route definitions

- 📄 **PostgreSQL Integration**
  - Database schema with proper foreign keys
  - Migration-ready schema setup

- 🐳 **(Optional) Docker Support**
  - Dockerfile and docker-compose (planned/coming soon)

---

## 📁 Project Structure

```

RestaurantManagementSystemAPI/
├── handlers/             # All route handlers
├── dbHelpers/            # DB interaction logic
├── middleware/           # JWT, sessions, RBAC
├── models/               # Struct definitions
├── utils/                # Helper utilities (e.g., token generation)
├── server.go             # Router and server initialization
├── go.mod / go.sum       # Dependencies
├── .env                  # Config (DB creds, JWT secrets)
└── README.md             # This file

````

---

## 🔧 Setup Instructions

### 1. Clone the Repo

```bash
git clone https://github.com/hemant-giri2004/RestaurantManagementSystemAPI.git
cd RestaurantManagementSystemAPI
````

### 2. Setup `.env` File

Create a `.env` file and add the following:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=restaurant_db
JWT_SECRET=your_jwt_secret
REFRESH_SECRET=your_refresh_secret
```

### 3. Run PostgreSQL

Ensure PostgreSQL is running and a database named `restaurant_db` is created.

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Server

```bash
go run server.go
```

---

## 🧪 API Endpoints Overview

### 🔐 Auth

| Method | Endpoint         | Description                 |
| ------ | ---------------- | --------------------------- |
| POST   | `/register`      | Register user with roles    |
| POST   | `/login`         | Login and get JWTs          |
| POST   | `/refresh-token` | Refresh JWT using session   |
| POST   | `/logout`        | Logout and invalidate token |

### 👥 Users

| Method | Endpoint    | Description              |
| ------ | ----------- | ------------------------ |
| GET    | `/users/me` | Get current user profile |

### 🏪 Restaurants

| Method | Endpoint           | Role Required   |
| ------ | ------------------ | --------------- |
| POST   | `/restaurants`     | Admin/Sub-Admin |
| GET    | `/restaurants`     | All Roles       |
| PUT    | `/restaurants/:id` | Admin/Sub-Admin |
| DELETE | `/restaurants/:id` | Admin/Sub-Admin |

### 🍽️ Dishes

| Method | Endpoint                  | Role Required   |
| ------ | ------------------------- | --------------- |
| POST   | `/restaurants/:id/dishes` | Admin/Sub-Admin |
| GET    | `/restaurants/:id/dishes` | All Roles       |
| PUT    | `/dishes/:id`             | Admin/Sub-Admin |
| DELETE | `/dishes/:id`             | Admin/Sub-Admin |

---

## 🛡️ Security & Best Practices

* Hashed passwords using `bcrypt`
* Environment-based config loading
* Token expiry and session tracking
* Middleware for verifying roles and tokens

---

## 📌 Roadmap / To-Do

* [ ] Add Docker and docker-compose support
* [ ] Add Swagger/OpenAPI documentation
* [ ] Implement pagination for listing endpoints
* [ ] Add unit tests and integration tests
* [ ] Rate limiting / request throttling

---

## 🤝 Contributing

Pull requests are welcome! Feel free to fork this repo and suggest improvements.

---


## 📬 Contact

**Hemant Giri**
📧 Email: [hemantgiri2004@gmail.com](mailto:hemantgiri2004@gmail.com)
🔗 GitHub: [@hemant-giri2004](https://github.com/hemant-giri2004)

---

```

---

Let me know if you want a **light/dark theme badge section**, **Swagger UI support**, or **GIF/screenshots** included later when you host it. I can help format them too.

Would you like me to push this README to your repo with a PR or just save it as a `.md` file for you?
```
