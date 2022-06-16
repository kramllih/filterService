package controllers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/kramllih/filterService/internal/database"
	"github.com/sirupsen/logrus"
)

const languageservice string = "/api/banned"

var (
	lvl1Heading = regexp.MustCompile(`(?:^|\s)(?:[#]\ )`)
	links       = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
)

func (c *Controller) Validate(ctx *gin.Context) {

	var message database.Message

	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	rdr := strings.NewReader(message.Body)

	scanner := bufio.NewScanner(rdr)
	scanner.Split(bufio.ScanLines)

	var txtlines []string

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		txtlines = append(txtlines, strings.TrimSpace(text))
	}

	if !lvl1Heading.MatchString(txtlines[0]) {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("first line must be a level 1 heading"))
		return
	}

	if len(txtlines) == 1 {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("body must contain at least 1 paragraph of text"))
		return
	}

	mes, _ := c.DB.GetMessage(message.ID)

	if mes != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("id already exists"))
		return
	}

	rejected, approvalRequired, err := c.handleValidation(&message, txtlines)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if approvalRequired {
		c.log.WithField("messageId", message.ID).Infof("message with ID [%s] requires approval", message.ID)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "your message is awaiting approval as it contains image links.",
		})
		return
	}

	if rejected {
		c.log.WithFields(logrus.Fields{"messageId": message.ID, "reason": message.Reason}).Infof("message with ID [%s] has been rejected.", message.ID)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "your message has has been rejected.",
			"reason": message.Reason,
		})
		return
	}

	c.log.WithField("messageId", message.ID).Infof("message with ID [%s] has been validated.", message.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "your message has been stored.",
	})

}

func (c *Controller) handleValidation(message *database.Message, txtlines []string) (bool, bool, error) {

	if message.Status == "" {
		message.Status = "pending"

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			return false, false, err
		}

		if err := c.DB.StoreMessage(message.ID, jsonMessage); err != nil {
			return false, false, errors.New("unable to store message")
		}
	}

	//handle reprocessed messages
	if len(message.Actions) > 0 {
		approved := 0

		for _, action := range message.Actions {
			if action.Status == "approved" {
				approved++
			}
		}

		if len(message.Actions) == approved {
			message.Status = "validated"

			jsonMessage, err := json.Marshal(message)
			if err != nil {
				return false, false, err
			}

			if err := c.DB.UpdateMessage(message.ID, jsonMessage); err != nil {
				return false, false, errors.New("unable to store message")
			}

			return false, false, nil

		} else {
			return false, false, nil
		}
	}

	actions := []database.Action{}
	approvalRequired := false
	rejected := false

	message.Status = "validated"

	banned, err := c.getBannedWords()
	if err != nil {
		return false, false, err
	}

	for _, line := range txtlines[1:] {

		//checking banned words
		matchedWords, err := c.checkAndHandleBannedWords(line, banned)
		if err != nil {
			return false, false, err
		}

		if matchedWords != nil {

			rejected = true

			message.Status = "rejected"
			message.Reason = fmt.Sprintf("message body contains these banned words: [%v]", strings.Join(matchedWords, ","))

			jsonMessage, err := json.Marshal(message)
			if err != nil {
				return false, false, err
			}

			if err := c.DB.StoreReject(message.ID, jsonMessage); err != nil {
				return false, false, errors.New("unable to store rejected message")
			}
			break

		}

		//checking links
		if links.MatchString(line) {
			act, required, reject, err := c.checkAndHandleLinks(line, message.ID)
			if err != nil {
				return false, false, err
			}

			if reject {
				rejected = true

				message.Status = "rejected"
				message.Reason = "message body contains external links"

				jsonMessage, err := json.Marshal(message)
				if err != nil {
					return false, false, err
				}

				if err := c.DB.StoreReject(message.ID, jsonMessage); err != nil {
					return false, false, errors.New("unable to store rejected message")
				}
			}

			if required {
				approvalRequired = true
			}

			if act.ID != "" {
				actions = append(actions, act)
			}
		}

	}

	message.Actions = actions

	if approvalRequired {
		message.Status = "awaiting approval"
		message.Reason = "message contains image that require approval"
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return false, false, err
	}

	if err := c.DB.UpdateMessage(message.ID, jsonMessage); err != nil {
		return false, false, errors.New("unable to store message")
	}

	return rejected, approvalRequired, nil

}

func (c *Controller) Rejected(ctx *gin.Context) {

	rejected, err := c.DB.GetAllRejected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	now := time.Now().UTC()

	ctx.JSON(http.StatusOK, gin.H{
		"updated":  now,
		"rejected": rejected,
	})

}

func (c *Controller) AllMessages(ctx *gin.Context) {

	approvals, err := c.DB.GetAllMessages()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	now := time.Now().UTC()

	ctx.JSON(http.StatusOK, gin.H{
		"updated":  now,
		"messages": approvals,
	})
}

func (c *Controller) checkAndHandleLinks(line, messageID string) (database.Action, bool, bool, error) {

	matches := links.FindStringSubmatch(line)
	if strings.HasPrefix(strings.TrimSpace(matches[2]), "http") {

		ok, _ := isImage(strings.TrimSpace(matches[2]))
		if !ok {
			return database.Action{}, false, true, nil
		}

		id, err := uuid.NewV4()
		if err != nil {
			return database.Action{}, false, false, fmt.Errorf("error generating ID: %w", err)
		}

		act := database.Action{
			ID:     id.String(),
			Status: "pending",
			Reason: fmt.Sprintf("image [%s] requires approval", matches[2]),
		}

		approval := database.Approval{
			ID:        id.String(),
			Status:    "pending",
			MessageID: messageID,
			Reason:    fmt.Sprintf("image [%s] requires approval", matches[2]),
		}

		jsonApproval, err := json.Marshal(approval)
		if err != nil {

			return database.Action{}, false, false, fmt.Errorf("error encoding json: %w", err)
		}

		if err := c.DB.StoreApproval(approval.ID, jsonApproval); err != nil {

			return database.Action{}, false, false, fmt.Errorf("unable to store message: %w", err)
		}

		return act, true, false, nil

	}

	return database.Action{}, false, false, nil
}

func isImage(url string) (bool, error) {

	res, err := http.Head(url)
	if err != nil {
		return false, fmt.Errorf("error access linked url: %w", err)
	}
	contentType := res.Header["Content-Type"]
	if strings.HasPrefix(contentType[0], "image") {
		return true, nil
	}

	return false, nil
}

func (c *Controller) getBannedWords() ([]string, error) {

	uri := c.httpClient.GetURI()
	if !strings.HasSuffix(uri, "banned") {
		c.httpClient.SetURI(c.httpClient.GetURI() + languageservice)
	}

	res, err := c.httpClient.FetchContent()
	if err != nil {
		return nil, err
	}

	workResp := struct {
		Updated time.Time
		Words   []string
	}{}

	if err := json.Unmarshal(res, &workResp); err != nil {
		return nil, err
	}

	return workResp.Words, nil
}

func (c *Controller) checkAndHandleBannedWords(line string, banned []string) ([]string, error) {

	var matchedWords []string

	words := strings.Fields(line)

	for _, word := range words {
		for _, bannedWord := range banned {
			//if test := strings.Index(strings.ToLower(word), bannedWord); test > -1 { //used for sub matching rather than matching the whole word
			if test := strings.EqualFold(strings.ToLower(word), bannedWord); test {
				matchedWords = append(matchedWords, word)
			}
		}
	}

	if len(matchedWords) > 0 {
		return matchedWords, nil
	}

	return nil, nil

}
