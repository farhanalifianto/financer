# 💰 Fintrac — Microservices Budget & Transaction Manager

Financer is a microservices-based financial tracking application built using **Go**, **gRPC**, **Docker**, **Redis**, and **PostgreSQL**.  
It allows users to manage **income/expense categories**, set budgets, create transactions, and calculate balances per category — all through a clean and modular microservices architecture.

---

## 🚀 Features

- 🔐 **User Authentication** (JWT-based)  
- 🧾 **Category Service** — Create and manage income/expense categories with budget limits.  
- 💳 **Transaction Service** — Record user transactions, calculate balances per category, and enforce category ownership.  
- 🔗 **gRPC Communication** between services.  
- 🐳 **Dockerized Microservices** for easy deployment.  
- ⚡ **Redis Caching** for faster category lookups.  
- 🌐 **Nginx Reverse Proxy** as the API gateway.

---

## 🧱 Project Structure

```text
fintrac-main/
├── category-service/
│ ├── controller/
│ ├── grpc_server/
│ ├── grpc_client/
│ ├── model/
│ ├── proto/
│ └── ...
│
├── transaction-service/
│ ├── controller/
│ ├── grpc_client/
│ ├── model/
│ ├── proto/
│ └── ...
│
├── user-service/
│ ├── controller/
│ ├── grpc_server/
│ ├── model/
│ ├── proto/
│ └── ...
│
├── docker-compose.yml
└── README.md

```
---

## 🛠️ Tech Stack

| Component              | Technology               |
|-------------------------|---------------------------|
| Language               | Go (Golang)              |
| Database               | PostgreSQL               |
| Messaging              | gRPC                    |
| Cache                   | Redis                   |
| Containerization       | Docker, Docker Compose  |
| API Framework          | Fiber                   |

---

## ⚡ Getting Started

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/farhanalifianto/financer.git
cd financer-main
```
### 2️⃣ Clone the Repository
```bash
docker-compose up --build
```
This will start:

- Category service (REST + gRPC Server and client)

- Transaction service (REST + gRPC client)

- User service (REST + gRPC Server)

- PostgreSQL containers for each service

- Redis

- Nginx reverse proxy
  
### 4️⃣ Access the Services

get postman json file in root folder

📝 API Overview

User Endpoints

- POST /user/register → Create a User

- POST /user/login → Login

Category Endpoints

- POST /category → Create a category (type: income or expense, with optional budget)

- GET /category/:id → Get category info (used by gRPC)

Transaction Endpoints

- POST /transaction → Create transaction (only using user-owned category)

- GET /transaction/balance → Get total balance

- GET /transaction/budget → Get budget usage and status

- GET /transaction/:id → Get transaction id

- GET /transaction/budget → Get budget usage per category

## 🤝 Contributing

Pull requests and suggestions are welcome!
If you'd like to contribute:

Fork the repo

Create a feature branch

Commit your changes

Open a Pull Request 🚀

## 🧠 Future Improvements

✅ Centralized Auth Service (JWT)

📈 Dashboard for analytics

🌐 Frontend UI (React / Next.js)

🚀 Kubernetes deployment

🙌 Acknowledgements

Built with ❤️ using Go and gRPC.



## Contact

Farhan Alifianto - [https://github.com/farhanalifianto/fintrac](https://github.com/farhanalifianto/fintrac) - farhanalifianto@gmail.com





