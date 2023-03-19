// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/deleteGroup": {
            "get": {
                "description": "删除群聊或好友",
                "tags": [
                    "通用"
                ],
                "summary": "删除群聊或好友",
                "parameters": [
                    {
                        "description": "删除群聊或好友入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.DeleteGroupParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/getBotList": {
            "get": {
                "description": "获取机器人qq号",
                "tags": [
                    "通用"
                ],
                "summary": "获取机器人qq号",
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/getFriendList": {
            "get": {
                "description": "获取好友列表",
                "tags": [
                    "通用"
                ],
                "summary": "获取好友列表",
                "parameters": [
                    {
                        "description": "获取好友列表入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.BotParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/getGroupList": {
            "get": {
                "description": "获取群列表",
                "tags": [
                    "通用"
                ],
                "summary": "获取群列表",
                "parameters": [
                    {
                        "description": "获取群列表入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.BotParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/getPermCode": {
            "get": {
                "description": "授权码",
                "tags": [
                    "用户"
                ],
                "summary": "授权码",
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "array",
                                            "items": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/getRequestList": {
            "get": {
                "description": "获取所有的请求",
                "tags": [
                    "通用"
                ],
                "summary": "获取所有的请求",
                "parameters": [
                    {
                        "description": "获取所有的请求入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.BotParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/types.RequestVo"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/getUserInfo": {
            "get": {
                "description": "获得用户信息",
                "tags": [
                    "用户"
                ],
                "summary": "获得用户信息",
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "$ref": "#/definitions/types.UserInfoVo"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/handleRequest": {
            "post": {
                "description": "处理一个请求",
                "tags": [
                    "通用"
                ],
                "summary": "处理一个请求",
                "parameters": [
                    {
                        "description": "处理一个请求入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.HandleRequestParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/job/add": {
            "post": {
                "description": "添加任务",
                "tags": [
                    "任务"
                ],
                "summary": "添加任务",
                "parameters": [
                    {
                        "description": "添加任务入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/job.Job"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/job/delete": {
            "post": {
                "description": "删除任务",
                "tags": [
                    "任务"
                ],
                "summary": "删除任务",
                "parameters": [
                    {
                        "description": "删除任务的入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/job.DeleteReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/job/list": {
            "get": {
                "description": "任务列表",
                "tags": [
                    "任务"
                ],
                "summary": "任务列表",
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "$ref": "#/definitions/job.JobListRsp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "前端登录",
                "tags": [
                    "用户"
                ],
                "summary": "前端登录",
                "parameters": [
                    {
                        "description": "前端登录",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.LoginParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "$ref": "#/definitions/types.LoginResultVo"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/logout": {
            "get": {
                "description": "登出",
                "tags": [
                    "用户"
                ],
                "summary": "登出",
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/manage/getAllPlugin": {
            "get": {
                "description": "获取所有插件的状态",
                "tags": [
                    "插件"
                ],
                "summary": "获取所有插件的状态",
                "parameters": [
                    {
                        "description": "获取所有插件的状态入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.AllPluginParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/types.PluginVo"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/manage/getPlugin": {
            "get": {
                "description": "获取某个插件的状态",
                "tags": [
                    "插件"
                ],
                "summary": "获取某个插件的状态",
                "parameters": [
                    {
                        "description": "获取某个插件的状态入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.PluginParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "$ref": "#/definitions/types.PluginVo"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/manage/updateAllPluginStatus": {
            "post": {
                "description": "更改某群所有插件状态",
                "tags": [
                    "响应"
                ],
                "summary": "更改某群所有插件状态",
                "parameters": [
                    {
                        "description": "修改插件状态入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.AllPluginStatusParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/manage/updatePluginStatus": {
            "post": {
                "description": "更改某一个插件状态",
                "tags": [
                    "插件"
                ],
                "summary": "更改某一个插件状态",
                "parameters": [
                    {
                        "description": "修改插件状态入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.PluginStatusParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/manage/updateResponseStatus": {
            "post": {
                "description": "更改某一个群响应",
                "tags": [
                    "响应"
                ],
                "summary": "更改某一个群响应",
                "parameters": [
                    {
                        "description": "修改群响应入参",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.ResponseStatusParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/api/sendMsg": {
            "post": {
                "description": "前端调用发送信息",
                "tags": [
                    "通用"
                ],
                "summary": "前端调用发送信息",
                "parameters": [
                    {
                        "description": "发消息参数",
                        "name": "object",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/types.SendMsgParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "result": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "job.DeleteReq": {
            "description": "删除任务的入参",
            "type": "object",
            "properties": {
                "idList": {
                    "description": "任务id",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "selfId": {
                    "description": "机器人qq",
                    "type": "integer"
                }
            }
        },
        "job.Job": {
            "description": "添加任务的入参",
            "type": "object",
            "properties": {
                "answerType": {
                    "description": "回答类型, jobType=3使用的参数, 1=文本消息, 2=注入消息",
                    "type": "integer"
                },
                "fullMatchType": {
                    "description": "指令别名类型, jobType=1使用的参数, 1=无状态消息, 2=主人消息",
                    "type": "integer"
                },
                "groupId": {
                    "description": "群聊id, jobType=2,3使用的参数, jobType=2且私聊, group_id=0",
                    "type": "integer"
                },
                "handler": {
                    "description": "执行内容",
                    "type": "string"
                },
                "id": {
                    "description": "任务id",
                    "type": "integer"
                },
                "jobType": {
                    "description": "任务类型,1-指令别名,2-定时任务,3-你问我答",
                    "type": "integer"
                },
                "matcher": {
                    "description": "当jobType=1时 为指令别名,当jobType=2时 为cron表达式,当jobType=3时 为正则表达式",
                    "type": "string"
                },
                "questionType": {
                    "description": "问题类型, jobType=3使用的参数, 1=单独问, 2=所有人问",
                    "type": "integer"
                },
                "selfId": {
                    "description": "机器人id",
                    "type": "integer"
                },
                "userId": {
                    "description": "用户id, jobType=2,3使用的参数, 当jobType=3, QuestionType=2,userId=0",
                    "type": "integer"
                }
            }
        },
        "job.JobListRsp": {
            "description": "任务列表的出参",
            "type": "object",
            "properties": {
                "jobList": {
                    "description": "任务列表",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/job.Job"
                    }
                }
            }
        },
        "types.AllPluginParams": {
            "description": "GetAllPlugin的入参",
            "type": "object",
            "properties": {
                "groupId": {
                    "description": "群id, gid\u003e0为群聊,gid\u003c0为私聊,gid=0为全部群聊",
                    "type": "integer"
                }
            }
        },
        "types.AllPluginStatusParams": {
            "description": "UpdateAllPluginStatus的入参",
            "type": "object",
            "properties": {
                "groupId": {
                    "description": "群id, gid\u003e0为群聊,gid\u003c0为私聊,gid=0为全部群聊",
                    "type": "integer"
                },
                "status": {
                    "description": "插件状态,0=禁用,1=启用,2=还原",
                    "type": "integer"
                }
            }
        },
        "types.BotParams": {
            "description": "GetGroupList,GetFriendList的入参",
            "type": "object",
            "properties": {
                "selfId": {
                    "description": "机器人qq",
                    "type": "integer"
                }
            }
        },
        "types.DeleteGroupParams": {
            "description": "退群或删除好友的入参",
            "type": "object",
            "properties": {
                "groupId": {
                    "description": "群id, gid\u003e0为群聊,gid\u003c0为私聊,gid=0为全部群聊",
                    "type": "integer"
                },
                "selfId": {
                    "description": "机器人qq",
                    "type": "integer"
                }
            }
        },
        "types.HandleRequestParams": {
            "description": "处理事件的入参",
            "type": "object",
            "properties": {
                "approve": {
                    "description": "是否同意, true=同意,false=拒绝",
                    "type": "boolean"
                },
                "flag": {
                    "description": "事件的flag",
                    "type": "string"
                },
                "reason": {
                    "description": "事件的原因, 拒绝的时候需要填",
                    "type": "string"
                }
            }
        },
        "types.LoginParams": {
            "description": "登录参数",
            "type": "object",
            "properties": {
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "types.LoginResultVo": {
            "description": "登录返回参数",
            "type": "object",
            "properties": {
                "desc": {
                    "description": "描述",
                    "type": "string"
                },
                "realName": {
                    "description": "实际名",
                    "type": "string"
                },
                "roles": {
                    "description": "角色",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.RoleInfo"
                    }
                },
                "token": {
                    "description": "token",
                    "type": "string"
                },
                "userId": {
                    "description": "用户id",
                    "type": "integer"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "types.PluginParams": {
            "description": "GetPlugin的入参",
            "type": "object",
            "properties": {
                "groupId": {
                    "description": "群id, gid\u003e0为群聊,gid\u003c0为私聊,gid=0为全部群聊",
                    "type": "integer"
                },
                "name": {
                    "description": "插件名",
                    "type": "string"
                }
            }
        },
        "types.PluginStatusParams": {
            "description": "UpdatePluginStatus的入参",
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "groupId": {
                    "description": "群id, gid\u003e0为群聊,gid\u003c0为私聊,gid=0为全部群聊",
                    "type": "integer"
                },
                "name": {
                    "description": "插件名",
                    "type": "string"
                },
                "status": {
                    "description": "插件状态,0=禁用,1=启用,2=还原",
                    "type": "integer"
                }
            }
        },
        "types.PluginVo": {
            "description": "全部插件的返回",
            "type": "object",
            "properties": {
                "banner": {
                    "description": "头像",
                    "type": "string"
                },
                "brief": {
                    "description": "简述",
                    "type": "string"
                },
                "id": {
                    "description": "插件序号",
                    "type": "integer"
                },
                "name": {
                    "description": "插件名",
                    "type": "string"
                },
                "pluginStatus": {
                    "description": "插件状态,false=禁用,true=启用",
                    "type": "boolean"
                },
                "responseStatus": {
                    "description": "响应状态, false=沉默,true=响应",
                    "type": "boolean"
                },
                "usage": {
                    "description": "用法",
                    "type": "string"
                }
            }
        },
        "types.RequestVo": {
            "description": "请求返回",
            "type": "object",
            "properties": {
                "comment": {
                    "description": "注释",
                    "type": "string"
                },
                "flag": {
                    "description": "请求flag",
                    "type": "string"
                },
                "groupId": {
                    "description": "群id",
                    "type": "integer"
                },
                "groupName": {
                    "description": "群名",
                    "type": "string"
                },
                "nickname": {
                    "description": "昵称",
                    "type": "string"
                },
                "requestType": {
                    "description": "请求类型",
                    "type": "string"
                },
                "selfId": {
                    "description": "机器人qq",
                    "type": "integer"
                },
                "subType": {
                    "description": "请求子类型",
                    "type": "string"
                },
                "userId": {
                    "description": "用户id",
                    "type": "integer"
                }
            }
        },
        "types.Response": {
            "description": "包装返回体",
            "type": "object",
            "properties": {
                "code": {
                    "description": "错误码",
                    "type": "integer"
                },
                "message": {
                    "description": "错误信息",
                    "type": "string"
                },
                "result": {
                    "description": "数据"
                },
                "type": {
                    "description": "待定",
                    "type": "string"
                }
            }
        },
        "types.ResponseStatusParams": {
            "description": "UpdateResponseStatus的入参",
            "type": "object",
            "properties": {
                "groupId": {
                    "description": "群id, gid\u003e0为群聊,gid\u003c0为私聊,gid=0为全部群聊",
                    "type": "integer"
                },
                "status": {
                    "description": "响应状态,0=沉默,1=响应",
                    "type": "integer"
                }
            }
        },
        "types.RoleInfo": {
            "description": "角色参数",
            "type": "object",
            "properties": {
                "roleName": {
                    "description": "角色名",
                    "type": "string"
                },
                "value": {
                    "description": "角色值",
                    "type": "string"
                }
            }
        },
        "types.SendMsgParams": {
            "description": "处理事件的入参",
            "type": "object",
            "properties": {
                "gidList": {
                    "description": "群聊数组",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "message": {
                    "description": "CQ码格式的消息",
                    "type": "string"
                },
                "selfId": {
                    "description": "机器人qq",
                    "type": "integer"
                }
            }
        },
        "types.UserInfoVo": {
            "description": "用户信息",
            "type": "object",
            "properties": {
                "avatar": {
                    "description": "头像",
                    "type": "string"
                },
                "desc": {
                    "description": "描述",
                    "type": "string"
                },
                "homePath": {
                    "description": "主页路径",
                    "type": "string"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "realName": {
                    "description": "实际名",
                    "type": "string"
                },
                "roles": {
                    "description": "角色",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.RoleInfo"
                    }
                },
                "token": {
                    "description": "token",
                    "type": "string"
                },
                "userId": {
                    "description": "用户id",
                    "type": "integer"
                },
                "username": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "zbp api",
	Description:      "zbp restful api document",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
