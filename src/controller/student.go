package controller

import (
	"net/http"
	"personality-teaching/src/code"
	"personality-teaching/src/logger"
	"personality-teaching/src/logic"
	"personality-teaching/src/model"
	"personality-teaching/src/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreateStudent(c *gin.Context) {
	var req model.CreateStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}

	resp, err := logic.NewStudentService().CreateStudent(req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("CreateStudent error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}

func AddStudentToClass(c *gin.Context) {
	var req model.AddStudentToClassReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	// 校验教师是否有该班级的权限和班级是否存在
	teacherID := c.GetString(utils.TeacherID)
	ok, err := logic.NewClassService().CheckPermission(teacherID, req.ClassID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("CheckPermission error: ", zap.Error(err))
		return
	}
	if !ok {
		code.CommonResp(c, http.StatusOK, code.InvalidPermission, code.EmptyData)
		return
	}
	// 将学生加入班级
	resp, err := logic.NewStudentService().UpdateClassID(req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("UpdateClassID error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}

func StudentsInClass(c *gin.Context) {
	var req model.ClassStudentListReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	// 校验教师是否有该班级的权限和班级是否存在
	teacherID := c.GetString(utils.TeacherID)
	ok, err := logic.NewClassService().CheckPermission(teacherID, req.ClassID)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("CheckPermission error: ", zap.Error(err))
		return
	}
	if !ok {
		code.CommonResp(c, http.StatusOK, code.InvalidPermission, code.EmptyData)
		return
	}
	resp, err := logic.NewStudentService().GetStudentsInClass(req)
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("GetStudentInClass error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}

// StudentNotInClass 查询未加入班级的学生
func StudentNotInClass(c *gin.Context) {
	var req model.EmptyClassStudentReq
	if err := c.ShouldBind(&req); err != nil {
		code.CommonResp(c, http.StatusBadRequest, code.InvalidParam, code.EmptyData)
		return
	}
	resp, err := logic.NewStudentService().GetStudentsInClass(model.ClassStudentListReq{
		ClassID:  utils.EmptyClassID,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		code.CommonResp(c, http.StatusInternalServerError, code.ServerBusy, code.EmptyData)
		logger.L.Error("GetStudentInClass error: ", zap.Error(err))
		return
	}
	code.CommonResp(c, http.StatusOK, code.Success, resp)
}
