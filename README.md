# REDDITCLONE
> A CRUD application that replicates the main functionality of Reddit

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation Instructions](#installation-instructions)
- [Usage](#usage)
- [Technologies Used](#technologies-used)

## Introduction
**RedditClone** is a backend-focused project written in **Go** (Golang) for VK Education homework. It replicates Reddit's core functionality through a custom API.  
The application leverages modern technologies such as Docker, Redis, MongoDB, and PostgreSQL to deliver a scalable and efficient service for managing and retrieving data.

## Features
- **Custom API**: Provides endpoints to interact with data (_users, cookies, posts, comments, likes_).  
- **Database Support**: Utilizes PostgreSQL for structured data (users) and MongoDB for unstructured storage (posts, comments, likes).  
- **User Cookies**: Implements user sessions using **JWT**, with cookies stored in Redis for fast access.  
- **Containerization**: Fully containerized using Docker for seamless deployment and scalability.

## Installation Instructions

### Prerequisites
Ensure you have the following installed:
- Docker
- Docker Compose

### Steps
1. Clone the Repository:
```bash
git clone https://github.com/myacey/redditclone.git
cd redditclone
```
2. Build and Run with Docker:
Use Docker Compose to build and run the application:
```bash
docker compose up --build
```
3. Accessing the API:
Once the application is running, you can access the frontend at `http://localhost:8080`.

## Usage

The API provides various endpoints to interact with the data. Below are some examples:

### Example API Requests

- **Register**: `POST /api/register` | _Sign up_
```bash
curl -X POST http://localhost:8080/api/register \  
 -H "Content-Type: application/json" \  
 -d '{"username":"your_username", "password":"your_password"}'
```

- **Login**: `POST /api/login` | _Sign in_
```bash
curl -X POST http://localhost:8080/api/login \  
 -H "Content-Type: application/json" \  
 -d '{"username":"your_username", "password":"your_password"}'
```

- **Get All Posts**: `GET /api/posts` | _List all posts_
```bash
curl -X GET http://localhost:8080/api/posts
```

- **Get Single Post**: `GET /api/post/<id>` | _Retrieve a specific post, including comments and votes_
```bash
curl -X GET http://localhost:8080/api/post/<id>
```

- **Add Comment**: `POST /api/post/<id>` | _Add a comment to a post_
```bash
curl -X POST http://localhost:8080/api/post/<id> \  
 -H "Authorization: Bearer your_token" \  
 -H "Content-Type: application/json" \  
 -d '{"comment": "Your comment here"}'
```

- **Upvote Post**: `GET /api/post/<id>/upvote` | _Vote on a post_
```bash
curl -X GET http://localhost:8080/api/post/<id>/upvote \  
 -H "Authorization: Bearer your_token"
```

## Technologies Used
- **Programming Language**: Go (Golang)  
- **Databases**: PostgreSQL, MongoDB  
- **Caching**: Redis  
- **Containerization**: Docker & Docker Compose
