GET http://localhost:8080/api/healthz

POST http://localhost:8080/admin/reset

POST http://localhost:8080/api/users
Content-Type: application/json
{
    "email": "hello@world.com"
}


