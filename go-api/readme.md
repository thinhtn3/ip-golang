# Backend Flow

## Sending chat flow
Client sends post request to /chat/session/:sessionId/message/ with message and role in body data, and accessToken (retrieved by SB from client) in auth headers. <br>
Access token goes through middleware and is validated before hitting SendMessage function in chatHandler. <br>
Handler strips sessionID and userID and calls chatService. Service verifies session ownerships and if valid, insert new message to DB before returning message struct back to client. <br>

## Creating chat flow
Client sends post request to /chat/create with question_id and access token in the body request. <br>
Access token is validated then request goes to chatHandler to extract question_id and user_id. <br>
Call createSession from NewChatService. createSession maps the arguments to the corresponding column name. Then upsert into chat_sessions table.<br>
Return chat session ID to handler which then returns to client with session ID and response 200.<br>

## Config
Load all env keys to config package and function to initialize Supabase.<br>
Database is init at main.go and is to be passed to services when needed.<br>
Dependency injection: connected to database when server starts, then passing that connection explicitly to the services that needs them.<br>
### Dependency Injection
Creating a dependency (supabase client).<br>
Inject depedency where needed e.g client, handler or service. <br>
Normally create a struct with the dependency as data, then create an instance using the constructor. Finally pass in dependency in receiver (member) function where the dependency (supabase client) is needed.<br>


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

