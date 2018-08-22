package router

import (
	"github.com/gin-gonic/gin"
)

//注册
func signUp(ctx *gin.Context) {

}

//登录系统
func signIn(ctx *gin.Context) {

}

//通过微信登录系统
func signWithWeiXin(ctx *gin.Context) {
	session.Values[SessionUid] = 1
}

//退出系统
func signOut(ctx *gin.Context) {

}

func bindWeiXin(ctx *gin.Context) {

}
