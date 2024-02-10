# GO-SHOP

This is a Go project with Gin, GORM and PostgreSQL. You can create your own store, publish your products. There is a cart system. For authorization I use token-based authentication system.

# REST API

The REST API to the example app is described below.

## Sign Up / Register

### Request

`POST /api/signup`

    curl -X POST -H "Content-Type: application/json" http://localhost:8080/api/signup          

### Body

    {
        "email":"your_email",
        "password:"your_password",
        "username":"your_username",
    }

### Response

    403 Status Forbidden
    {
        "message": "You are signed up! Now verify your email"
    }

## Log In

### Request

`POST /api/login`

    curl -X POST -H "Content-Type: application/json" http://localhost:8080/api/login    

### Body

    {
        "password:"your_password",
        "username":"your_username",
    }

### Response

    200 OK
    {
        "token":"user.Token",
		"tokenExpireAt":"user.TokenExpireAt",
		"refreshToken":"user.RefreshToken",
		"refreshTokenExpireAt":"user.RefreshTokenExpireAt"
    }

## Delete Account

### Request

`POST /api/delete-account`

    curl -X POST -H "Content-Type: application/json" http://localhost:8080/api/delete-account          

### Body

    {
        "password:"your_password",
        "username":"your_username",
    }

### Response

    200 OK
    {
        "message": "We sent you an email to delete your account!"
    }

## Create Store

### Request

`POST /api/create-store`

    curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer [your token]" http://localhost:8080/api/create-store          

### Body

    {
        "name:"your_store_name",
        "linkName":"your_store_name_link with regex "^[a-zA-Z0-9 ]*$"",
    }

### Response

    200 OK
    {
        "message": "store created!"
    }

## Create Product For Store

### Request

`POST /api/create-store-item`

    curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer [your token]" http://localhost:8080/api/create-store-item          

### Body

    {
        "name:"your_product_name",
        "price":"price_in_float_format",
        "itemCount":"product_count_int"
    }

### Response

    200 OK
    {
        "message": "Item created!"
    }

## Change Password

### Request

`POST /api/change-password`

    curl -X POST -H "Content-Type: application/json" http://localhost:8080/api/change-password          

### Body

    {
        "password:"your_new_password",
        "token":"token_that_sent_to_your_email_from /api/forget-password token will look like; 'msASda123...|change-password'",
    }

### Response

    200 OK
    {
        "message": "Your password has changed!"
    }

## Refresh Expired Token With RefreshToken

### Request

`GET /api/refresh-token`

    curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer [your refresh token]" http://localhost:8080/api/refresh-token          

### Response

        200 OK
        {
            "token":"user.Token",
		    "tokenExpireAt":"user.TokenExpireAt",
		    "refreshToken":"user.RefreshToken",
		    "refreshTokenExpireAt":"user.RefreshTokenExpireAt"
        }

## User Actions; verify-email - delete-account

### Request

`GET /api/user-action?token=[token_that_sent_to_your_email]`

    curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/user-action          

### Response verify-email

        200 OK
        {
            "message": "Your email is verified!"
        }
### Response delete-account

        200 OK
        {
            "message": "User deleted!"
        }

## Forget Password

### Request

`GET /api/forget-password?email=[your_email]`

    curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/forget-password

### Response

        200 OK
        {
            "message": "We sent you an email to reset your password!"
        }

## Get User Public Data

### Request

`GET /api/user/:username`

    curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/user/:username
    
### Response

        200 OK
        {
            "message": "User found!",
            "user": {
                "username":  user.Username,
                "createdAt": user.CreatedAt,
            }
        }

## Get Store Data

### Request

`GET /api/store/:storeLinkName`

    curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/store/:storeLinkName          

### Response

        200 OK
        {
            "message": "Store found!",
            "store": {
                "name":store.Name,
        		"linkName":store.LinkName,
        		"createdAt":store.CreatedAt,
        		"itemCount":store.ItemCount,
        		"userId":store.UserId,
            },
            "cartItemLength":"length_int",
            "items": [
                {
                    "id":uint
                	"storeId":int
                	"createdAt":time.Time
                	"name":string
                	"itemCount":int
                	"itemSellCount":int
                	"price":float64
                }
            ]
        }

## Get Cart Data

### Request

`GET /api/cart`

    curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer [your token]" http://localhost:8080/api/cart          

### Response

        200 OK
        {
            "items": map[string]int [],
            "totalPrice":"float",
            "itemLength":"int"
        }

## Add Product To Cart

### Request

`GET /api/add-to-cart?item=[item_id]`

    curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer [your token]" http://localhost:8080/api/add-to-cart  

### Response

        200 OK
        {
            "message": "Item added to your cart!"
        }        

## Delete Product From Cart

### Request

`GET /api/delete-from-cart?item=[item_id]`

    curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer [your token]" http://localhost:8080/api/delete-from-cart        

### Response

        200 OK
        {
            "message": "Item deleted from your cart!"
        }       


