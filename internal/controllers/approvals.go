package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kramllih/filterService/internal/database"
)

func (c *Controller) Approve(ctx *gin.Context) {

	id := ctx.Param("id")

	approval, err := c.DB.GetApproval(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := c.DB.DeleteApprovals(approval.ID); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	message, err := c.DB.GetMessage(approval.MessageID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	actions := []database.Action{}
	approvedCount := 0

	for _, act := range message.Actions {
		if act.ID != approval.ID {
			if act.Status == "approved" {
				approvedCount++
			}
			actions = append(actions, act)
			continue
		}

		act.Status = "approved"
		actions = append(actions, act)
		approvedCount++
	}

	message.Actions = actions

	if len(actions) == approvedCount {
		_, _, err := c.handleValidation(message, nil)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	updatedMessage, err := json.Marshal(message)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := c.DB.UpdateMessage(approval.MessageID, updatedMessage); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}

func (c *Controller) Reject(ctx *gin.Context) {

	id := ctx.Param("id")

	approval, err := c.DB.GetApproval(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := c.DB.DeleteApprovals(approval.ID); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	message, err := c.DB.GetMessage(approval.MessageID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	actions := []database.Action{}
	rejectedCount := 0

	for _, act := range message.Actions {
		if act.ID == approval.ID {
			act.Status = "rejected"
			actions = append(actions, act)
			rejectedCount++
		}

	}

	if rejectedCount >= 1 {
		message.Status = "rejected"
	}

	message.Actions = actions

	updatedMessage, err := json.Marshal(message)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := c.DB.StoreReject(approval.MessageID, updatedMessage); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := c.DB.StoreMessage(approval.MessageID, updatedMessage); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}

func (c *Controller) AllApprovals(ctx *gin.Context) {

	approvals, err := c.DB.GetAllApprovals()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	now := time.Now().UTC()

	ctx.JSON(http.StatusOK, gin.H{
		"updated":   now,
		"approvals": approvals,
	})

}
