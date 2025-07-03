# 🍽️ Restaurant Management System API

Welcome to the **Restaurant Management System API** — a backend project built with 💙 **Golang** and 🐘 **PostgreSQL**. This system lets admins manage restaurants, dishes, and users securely using **JWT tokens** and **role-based access**.

Whether you're a developer, recruiter, or curious learner — this README will walk you through everything implemented.

---

## ✨ What This Project Does

This API system allows you to:

* 🔐 Register & log in users with secure authentication
* 🧑‍🍳 Add and manage restaurants
* 🍛 Add and manage dishes under those restaurants
* 🧑‍💼 Assign user roles like `Admin`, `Sub-Admin`, and `User`
* 🎯 Protect endpoints based on user roles
* 🧠 Use refresh tokens for session management
* 🧩 Maintain clean, modular project architecture

---

## 🔐 Authentication System

The system uses a powerful combo of:

* **JWT (JSON Web Tokens)** for access tokens
* **Refresh tokens** stored in the database for secure session renewal
* ✅ Supports login, logout, and token expiration handling

---

## 👥 User Roles & Access

Every user gets a role during registration:

| Role            | Can Do                                           |
| --------------- | ------------------------------------------------ |
| 👑 Admin        | Full control — manage restaurants, dishes, users |
| 🧑‍🔧 Sub-Admin | Can manage restaurants and dishes, but not users |
| 🙋‍♂️ User      | View their own profile and assigned data only    |

Middleware checks are applied behind the scenes to protect routes.

---

## 🍽️ Features at a Glance

### 🏢 Restaurant Management

* Add new restaurants
* Update or delete existing ones
* View list of all restaurants

### 🍛 Dish Management

* Add dishes to specific restaurants
* Update and delete them as needed
* Get list of all dishes under a restaurant

### 🙋‍♂️ User Management

* Register users with roles and address
* View personal profile after login

### 🧱 Clean Architecture

Everything is neatly separated:

* Route handlers 🚦
* Database query logic 📦
* Middleware 🔒
* Models & request/response structs 📑
* Utility helpers 🛠️

---

## 🧭 Folder Structure

Here’s a bird’s-eye view of the project:

```
RestaurantManagementSystemAPI/
├── handlers/         → All API route logic
├── dbHelpers/        → Raw SQL queries & DB operations
├── middleware/       → Auth & role check logic
├── models/           → Request & response data structs
├── utils/            → Token generation & helpers
├── server.go         → Entry point of the app
├── .env              → Environment config (DB creds, secrets)
├── go.mod / go.sum   → Go module dependencies
└── README.md         → Project documentation
```

---

## 💡 Why I Built This

I wanted to build a real-world backend system that mimics how companies manage restaurants and users securely — while learning about:

* ✅ Backend structuring in Go
* ✅ JWT authentication + sessions
* ✅ Database migrations and logic separation
* ✅ Role-based access like real SaaS apps

---

## 🚀 Future Plans

Here’s what I’m planning to add next:

* 📄 Swagger documentation for all APIs
* 🐳 Docker support for easy deployment
* 💻 Admin dashboard using React
* ✉️ Email verification during registration
* 🧪 Automated test coverage

---

## 🙋‍♂️ About Me

Hey! I’m **Hemant Giri** 👋
I’m an MCA student passionate about backend development, Golang, and building secure APIs.

* 🎓 GL Bajaj Institute of Technology and Management
* 🛠️ Strong in Go, SQL, C++, and system design
* 🎯 Goal: Become a backend developer in a top MNC

📍 From Bulandshahr, Uttar Pradesh

> Feel free to connect or explore my GitHub profile for more cool stuff.

---
