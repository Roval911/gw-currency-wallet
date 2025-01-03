// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "support@currencywallet.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/balance": {
            "get": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Retrieves the balance of the authenticated user's wallet",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallet"
                ],
                "summary": "Get wallet balance",
                "responses": {
                    "200": {
                        "description": "Balance retrieved",
                        "schema": {
                            "$ref": "#/definitions/storages.Wallet"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch balance",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/createwallet": {
            "post": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Creates a wallet for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallet"
                ],
                "summary": "Create a new wallet",
                "responses": {
                    "201": {
                        "description": "Wallet registered successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to create wallet",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/exchange": {
            "post": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Exchanges one currency to another in the authenticated user's wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Exchange"
                ],
                "summary": "Exchange currency",
                "parameters": [
                    {
                        "description": "Exchange details",
                        "name": "exchangeRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.ExchangeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Exchange successful",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to process exchange",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/exchange/rates": {
            "get": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Retrieves current exchange rates for supported currencies from the exchange service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Exchange"
                ],
                "summary": "Get exchange rates",
                "responses": {
                    "200": {
                        "description": "Exchange rates retrieved",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve exchange rates",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/login": {
            "post": {
                "description": "Authenticates user with username and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Log in a user",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "loginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Token returned",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid input data",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "Invalid username or password",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/register": {
            "post": {
                "description": "Creates a new user with provided credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User registered successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid input data",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to hash password",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/wallet/deposit": {
            "post": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Adds funds to the authenticated user's wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallet"
                ],
                "summary": "Deposit funds to wallet",
                "parameters": [
                    {
                        "description": "Deposit details",
                        "name": "depositRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.DepositRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Account topped up successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to deposit funds",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/v1/wallet/withdraw": {
            "post": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "Withdraws funds from the authenticated user's wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wallet"
                ],
                "summary": "Withdraw funds from wallet",
                "parameters": [
                    {
                        "description": "Withdrawal details",
                        "name": "withdrawRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.WithdrawRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Withdrawal successful",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to withdraw funds",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "storages.DepositRequest": {
            "type": "object",
            "required": [
                "amount",
                "currency"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "currency": {
                    "type": "string",
                    "enum": [
                        "USD",
                        "RUB",
                        "EUR"
                    ]
                }
            }
        },
        "storages.ExchangeRequest": {
            "type": "object",
            "required": [
                "amount",
                "from_currency",
                "to_currency"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "from_currency": {
                    "type": "string"
                },
                "to_currency": {
                    "type": "string"
                }
            }
        },
        "storages.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "password123"
                },
                "username": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "storages.User": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "username": {
                    "type": "string",
                    "minLength": 2
                }
            }
        },
        "storages.Wallet": {
            "type": "object",
            "properties": {
                "EUR": {
                    "type": "number"
                },
                "RUB": {
                    "type": "number"
                },
                "USD": {
                    "type": "number"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "storages.WithdrawRequest": {
            "type": "object",
            "required": [
                "amount",
                "currency"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "currency": {
                    "type": "string",
                    "enum": [
                        "USD",
                        "RUB",
                        "EUR"
                    ]
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Currency Wallet API",
	Description:      "API for managing wallets and currency exchanges.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
