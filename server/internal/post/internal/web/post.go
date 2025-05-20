package web

import (
	"fmt"

	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/post/internal/service"
	"github.com/gin-gonic/gin"
)

func NewPostHandler(serv service.IPostService) *PostHandler {
	return &PostHandler{
		serv: serv,
	}
}

type PostHandler struct {
	serv service.IPostService
}

func (h *PostHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin-api/post")
	{
		adminGroup.GET("draft/list", apiwrap.WrapWithBody(h.AdminGetDraftDetailPostList))
		adminGroup.GET("bin/list", apiwrap.WrapWithBody(h.AdminGetBinDetailPostList))
		adminGroup.POST("create", apiwrap.WrapWithBody(h.AdminCreatePost))
		adminGroup.PUT("update", apiwrap.WrapWithBody(h.AdminUpdatePost))
		adminGroup.PUT("update/publish-status", apiwrap.WrapWithBody(h.AdminUpdatePostPublishStatus))
		adminGroup.PUT("restore/:id", apiwrap.WrapWithUri(h.AdminRestorePost))
		adminGroup.PUT("restore/batch", apiwrap.WrapWithBody(h.AdminRestorePostBatch))
		adminGroup.DELETE("soft-delete/:id", apiwrap.WrapWithUri(h.AdminSoftDeletePost))
		adminGroup.DELETE("soft-delete/batch", apiwrap.WrapWithBody(h.AdminSoftDeletePostBatch))
		adminGroup.DELETE("delete/:id", apiwrap.WrapWithUri(h.AdminDeletePost))
		adminGroup.DELETE("delete/batch", apiwrap.WrapWithBody(h.AdminDeletePostBatch))
	}
	postGroup := engine.Group("/post")
	{
		postGroup.GET("/detail/list", apiwrap.WrapWithBody(h.GetPublishDetailPostList))
		postGroup.GET("/detail/:id", apiwrap.WrapWithUri(h.GetDetailPostById))
		postGroup.GET("/:id", apiwrap.WrapWithUri(h.GetPostById))
		postGroup.GET("/search", apiwrap.Wrap(h.GetPostByKeyWord))
		postGroup.GET("/sitemap", apiwrap.Wrap(h.GetSiteMap))
	}
}

func (h *PostHandler) AdminCreatePost(c *gin.Context, postReq PostRequest) *apiwrap.Response[any] {
	post := h.PostDTOToDomain(postReq)
	err := h.serv.AdminCreatePost(c, post)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("创建文章成功")
}

func (h *PostHandler) AdminUpdatePost(c *gin.Context, postUpdateReq PostUpdateRequest) *apiwrap.Response[any] {
	postUpdate := h.PostUpdateDTOToDomain(postUpdateReq)
	err := h.serv.AdminUpdatePost(c, postUpdate)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新文章成功")
}

func (h *PostHandler) AdminUpdatePostPublishStatus(c *gin.Context, postPublishStatusRequest PostPublishStatusRequest) *apiwrap.Response[any] {
	err := h.serv.AdminUpdatePostPublishStatus(c, apiwrap.ConvertBsonID(postPublishStatusRequest.ID), *postPublishStatusRequest.IsPublish)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("更新文章发布状态成功")
}

func (h *PostHandler) AdminRestorePost(c *gin.Context, postIDRequest PostIDRequest) *apiwrap.Response[any] {
	err := h.serv.AdminRestorePost(c, apiwrap.ConvertBsonID(postIDRequest.ID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("恢复文章成功")
}

func (h *PostHandler) AdminRestorePostBatch(c *gin.Context, postIDListRequest PostIDListRequest) *apiwrap.Response[any] {
	fmt.Println("postIDListRequest", postIDListRequest)
	err := h.serv.AdminRestorePostBatch(c, apiwrap.ConvertBsonIDList(postIDListRequest.IDList))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("批量恢复文章成功")
}

func (h *PostHandler) AdminSoftDeletePost(c *gin.Context, postIDRequest PostIDRequest) *apiwrap.Response[any] {
	err := h.serv.AdminSoftDeletePost(c, apiwrap.ConvertBsonID(postIDRequest.ID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("软删除文章成功")
}

func (h *PostHandler) AdminSoftDeletePostBatch(c *gin.Context, postIDListRequest PostIDListRequest) *apiwrap.Response[any] {
	err := h.serv.AdminSoftDeletePostBatch(c, apiwrap.ConvertBsonIDList(postIDListRequest.IDList))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("批量软删除文章成功")
}

func (h *PostHandler) AdminDeletePost(c *gin.Context, postIDRequest PostIDRequest) *apiwrap.Response[any] {
	err := h.serv.AdminDeletePost(c, apiwrap.ConvertBsonID(postIDRequest.ID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("删除文章成功")
}

func (h *PostHandler) AdminDeletePostBatch(c *gin.Context, postIDListRequest PostIDListRequest) *apiwrap.Response[any] {
	err := h.serv.AdminDeletePostBatch(c, apiwrap.ConvertBsonIDList(postIDListRequest.IDList))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("批量删除文章成功")
}

func (h *PostHandler) GetPublishDetailPostList(c *gin.Context, pageReq apiwrap.Page) *apiwrap.Response[any] {
	postDetailList, total, err := h.serv.GetPostDetailList(c, &pageReq, "publish")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	postVos := h.PostDetailListToVOList(postDetailList)

	pageVo := apiwrap.ToPageVO(pageReq.PageNo, pageReq.PageSize, total, postVos)
	return apiwrap.SuccessWithDetail[any](pageVo, "获取发布文章详情列表成功")
}

func (h *PostHandler) AdminGetDraftDetailPostList(c *gin.Context, pageReq apiwrap.Page) *apiwrap.Response[any] {
	postDetailList, total, err := h.serv.GetPostDetailList(c, &pageReq, "draft")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	postVos := h.PostDetailListToVOList(postDetailList)

	pageVo := apiwrap.ToPageVO(pageReq.PageNo, pageReq.PageSize, total, postVos)
	return apiwrap.SuccessWithDetail[any](pageVo, "获取草稿箱文章详情列表成功")
}

func (h *PostHandler) AdminGetBinDetailPostList(c *gin.Context, pageReq apiwrap.Page) *apiwrap.Response[any] {
	postDetailList, total, err := h.serv.GetPostDetailList(c, &pageReq, "bin")
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	postVos := h.PostDetailListToVOList(postDetailList)

	pageVo := apiwrap.ToPageVO(pageReq.PageNo, pageReq.PageSize, total, postVos)
	return apiwrap.SuccessWithDetail[any](pageVo, "获取回收站文章详情列表成功")
}

func (h *PostHandler) GetDetailPostById(c *gin.Context, postIDRequest PostIDRequest) *apiwrap.Response[any] {
	postDetail, err := h.serv.GetPostDetailById(c, apiwrap.ConvertBsonID(postIDRequest.ID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.PostDetailToVO(postDetail), "获取文章详情成功")
}

func (h *PostHandler) GetPostById(c *gin.Context, postIDRequest PostIDRequest) *apiwrap.Response[any] {
	post, err := h.serv.GetPostById(c, apiwrap.ConvertBsonID(postIDRequest.ID))
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.PostToVO(post), "获取文章成功")
}

func (h *PostHandler) GetPostByKeyWord(c *gin.Context) *apiwrap.Response[any] {
	keyword := c.Query("keyword")
	postList, err := h.serv.GetPostByKeyWord(c, keyword)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](h.PostListToVOList(postList), "获取文章成功")
}

func (h *PostHandler) GetSiteMap(c *gin.Context) *apiwrap.Response[any] {
	postList, err := h.serv.GetAllPublishPost(c)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	siteMapVos := h.PostListToSiteMapVOList(postList)
	return apiwrap.SuccessWithDetail[any](siteMapVos, "获取站点地图成功")
}
