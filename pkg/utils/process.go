package utils

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	log "github.com/sirupsen/logrus"
	"regexp"
	"stock_scraper/types"
	"strings"
)

func diff(previousContent string, newContent string) ([]string, []string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(previousContent, newContent, false)
	var addedTexts []string
	var removedTexts []string
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			addedTexts = append(addedTexts, diff.Text)
		case diffmatchpatch.DiffDelete:
			removedTexts = append(removedTexts, diff.Text)
		}
	}

	return addedTexts, removedTexts
}

func applyRule(rule types.Rule, text string) *types.Action {
	switch rule.Strategy {
	case "has":
		if strings.Contains(
			strings.ToLower(text),
			strings.ToLower(rule.Text),
		) {
			return &types.Action{Type: "found", Content: text}
		}
	case "match":
		r, err := regexp.Compile("(?i)" + rule.Text)
		if err != nil {
			log.WithFields(log.Fields{
				"error": fmt.Sprintf("%s", err),
				"regex": rule.Text,
			}).Error("Invalid Regex")
		}
		if r.MatchString(text) {
			return &types.Action{Type: "found", Content: text}
		}
	}
	return nil
}

func ApplyRules(rules []types.Rule, previousContent string, newContent string) []types.Action {
	addedTexts, removedTexts := diff(previousContent, newContent)
	addedText := strings.Join(addedTexts, "")
	removedText := strings.Join(removedTexts, "")
	var actions []types.Action

	for _, rule := range rules {
		var action *types.Action
		switch rule.Condition {
		case "added":
			action = applyRule(rule, addedText)
		case "removed":
			action = applyRule(rule, removedText)
		case "text":
			action = applyRule(rule, newContent)
		}
		if action != nil {
			actions = append(actions, *action)
		}
	}

	return actions
}
