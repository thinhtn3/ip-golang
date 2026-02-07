# Backend Architecture Overview

This document summarizes the backend principles used in this project.
The goal is to build a **maintainable, testable, and scalable Go backend** using industry-standard patterns.


## Layered Architecture
handlers -> services -> domain (interfaces) <- infra (implementations) <br>

## 1. Handlers (Delivery Layer)
**Purpose**
- Handles HTTP requests and response between server and client
- Translating HTTP to GoLang domain/service types and call services.

**Responsibilities**
- Parse request data from client
- Extract user context and concert to GoLang domain types
- Call service mthods
- Convert response from domain to HTTP responses (return error if exist)

**How it works**
- Use gin.Context (wrapper around HTTP req/res) to parse into Go types.
