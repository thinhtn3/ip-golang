# Backend Architecture Overview

This document summarizes the backend principles used in this project.
The goal is to build a **maintainable, testable, and scalable Go backend** using industry-standard patterns.

<br>

## Layered Architecture
handlers -> services -> domain (interfaces) <- infra (implementations)

