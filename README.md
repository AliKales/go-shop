# GO-SHOP

This is a Go project with Gin, GORM and PostgreSQL. You can create your own store, publish your products. There is a cart system. For authorization I use token-based authentication system.

# REST API

The REST API to the example app is described below.

## Get list of Things

### Request

`GET /api/refresh-token`

    curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer [your refresh token]" http://localhost:8080/api/refresh-token

### Response

    {
        
    }
