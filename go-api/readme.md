# Backend Flow

## Config
Load all env keys to config package and function to initialize Supabase

## Middlewares
### auth.go (Required for all protected routes)
Receives access token from frontend via Auth Bearer Token, verify that token is from Supabase and is untampered with.
If validated, extract user info and pass onto next handler.
Must be called every time a request is made to a protected route as each context is a new request context.
Validates existing access token from client before making requests.
