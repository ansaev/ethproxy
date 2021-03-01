// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2021-03-01 15:26:32.830354116 +0300 MSK m=+3.516493046

package api

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/block/{blockID}/tx/{txID}": {
            "get": {
                "description": "Get tx by hash or index number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "GetTx handler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Block number, or 'latest'",
                        "name": "blockID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Tx index, or tx hash",
                        "name": "txID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.TxResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.TxResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.TxResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ErrorForm": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "api.TxResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/domain.Transaction"
                },
                "error": {
                    "$ref": "#/definitions/api.ErrorForm"
                },
                "ok": {
                    "type": "boolean"
                }
            }
        },
        "domain.Transaction": {
            "type": "object",
            "properties": {
                "blockHash": {
                    "type": "string"
                },
                "blockNumber": {
                    "type": "string"
                },
                "from": {
                    "type": "string"
                },
                "gas": {
                    "type": "string"
                },
                "gasPrice": {
                    "type": "string"
                },
                "hash": {
                    "type": "string"
                },
                "input": {
                    "type": "string"
                },
                "nonce": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                },
                "transactionIndex": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}