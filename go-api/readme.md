# Backend Flow

## Creating chat flow
Client sends post request to /chat/create with question_id and access token in the body request. <br>
Access token is validated then request goes to chatHandler to extract question_id and user_id. <br>
Call createSession from NewChatService. createSession maps the arguments to the corresponding column name. Then upsert into chat_sessions table.

## Config
Load all env keys to config package and function to initialize Supabase.<br>
Database is init at main.go and is to be passed to services when needed.<br>
Dependency injection: connected to database when server starts, then passing that connection explicitly to the services that needs them.<br>


## Middlewares
### auth.go (Required for all protected routes)
Receives access token from frontend via Auth Bearer Token, verify that token is from Supabase and is untampered with.
If validated, extract user info and pass onto next handler.
Must be called every time a request is made to a protected route as each context is a new request context.
Validates existing access token from client before making requests.

## Models
### chat_sessions
user_id column: foreign key that checks auth.users table to see if user_id exists.<br>
question_id column: foreign key which checks against question_bank table to see if question exist.<br>

### messages
user_id column: foreign key that checks against auth.users table.<br>
chat_session_id column: uuid column that is a foreign key and checks against chat sessions table. To be used to gather all messages and render.<br>

