package process

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	log "github.com/sirupsen/logrus"
	"regexp"
	"stock_scraper/internal/utils"
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

func applyRuleToContent(rule types.Rule, previousContent string, newContent string) *types.Action {
	addedTexts, removedTexts := diff(previousContent, newContent)
	addedText := strings.Join(utils.Map(addedTexts, strings.TrimSpace), "")
	removedText := strings.Join(utils.Map(removedTexts, strings.TrimSpace), "")

	var action *types.Action
	switch rule.Condition {
	case "changed":
		if previousContent != newContent {
			action = &types.Action{Type: "found", Content: newContent}
		}
	case "added":
		action = applyRule(rule, addedText)
	case "removed":
		action = applyRule(rule, removedText)
	case "text":
		action = applyRule(rule, newContent)
	}
	if action != nil {
		action.Diff = types.Diff{AddedText: addedText, RemovedText: removedText}
	}
	return action
}

func applyRulesToContent(item types.Item, rules []types.Rule, previousContent string, newContent string, result *types.ParsedResults, actionsMap *map[string]types.Action) {
	aMap := *actionsMap

	for _, rule := range rules {
		action := applyRuleToContent(rule, previousContent, newContent)
		if action != nil {
			action.Link = item.TrackedUrl
			if result != nil {
				if result.ItemLink != nil {
					action.Link = *result.ItemLink
				}
				if result.ItemAddToCartLink != nil {
					action.AddToCartLink = *result.ItemAddToCartLink
				}
			}
			aMap[action.Link] = *action
		}
	}
}

func ApplyRules(item types.Item, rules []types.Rule, previousResult *types.Result, newResult types.Result) []types.Action {
	if previousResult == nil {
		return []types.Action{}
	}

	actionsMap := make(map[string]types.Action)
	if len(previousResult.Results) != len(newResult.Results) ||
		len(previousResult.Results) < 2 ||
		len(newResult.Results) < 2 {
		previousContent := strings.TrimSpace(previousResult.Content)
		newContent := strings.TrimSpace(newResult.Content)
		applyRulesToContent(item, rules, previousContent, newContent, nil, &actionsMap)
	} else {
		for i, result := range newResult.Results {
			previousContent := strings.TrimSpace(previousResult.Results[i].Content)
			newContent := strings.TrimSpace(result.Content)

			applyRulesToContent(item, rules, previousContent, newContent, &result, &actionsMap)
		}
	}

	var actions []types.Action
	for _, value := range actionsMap {
		actions = append(actions, value)
	}

	return actions
}
