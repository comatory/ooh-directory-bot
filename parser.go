package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

type DomNodeNotFoundError struct{}

func (*DomNodeNotFoundError) Error() string {
	return "Required HTML structure not present"
}

type Result struct {
	url        string
	title      string
	summary    string
	authorName string
}

func newEmptyResult() Result {
	return Result{
		url:        "",
		title:      "",
		summary:    "",
		authorName: "",
	}
}

func hasClass(attributes *[]html.Attribute, name string) bool {
	for _, attr := range *attributes {
		if attr.Key == "class" && attr.Val == name {
			return true
		}
	}

	return false
}

func isWebsitesList(node *html.Node) bool {
	if node.Data != "ol" {
		return false
	}

	return hasClass(&node.Attr, "websites")
}

func isWebsiteListItem(node *html.Node) bool {
	if node.Data != "li" {
		return false
	}

	return hasClass(&node.Attr, "websites__item")
}

func isWebsiteDetail(node *html.Node) bool {
	if node.Data != "details" {
		return false
	}

	return hasClass(&node.Attr, "website__details")
}

func isWebsiteDetailBody(node *html.Node) bool {
	if node.Data != "div" {
		return false
	}

	return hasClass(&node.Attr, "website__details__body")
}

func findList(node *html.Node) (*html.Node, bool) {
	var listNode *html.Node
	var isFound = false

	var traverse func(node *html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && isWebsitesList(node) {
			isFound = true
			listNode = node
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}

	}

	traverse(node)

	return listNode, isFound
}

func getListItems(node *html.Node) []*html.Node {
	items := []*html.Node{}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isWebsiteListItem(child) {
			items = append(items, child)
		}
	}

	return items
}

func getLinkWithAnchor(node *html.Node) (string, string) {
	var url string
	title := fmt.Sprint(node.FirstChild)
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			url = attr.Val
		}
	}

	return url, title
}

func parseListItemDetails(node *html.Node) (Result, bool) {
	result := newEmptyResult()
	isOk := false

	var traverseDetailsBody func(node *html.Node, r *Result)
	traverseDetailsBody = func(node *html.Node, r *Result) {
		if node.Type == html.ElementNode && node.Data == "a" {
			url, title := getLinkWithAnchor(node)

			r.title = title
			r.url = url
		}

		if node.Type == html.ElementNode && node.Data == "small" {
			r.authorName = fmt.Sprint(node.FirstChild)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseDetailsBody(child, &result)
		}
	}

	var traverseDetails func(node *html.Node, r *Result)
	traverseDetails = func(node *html.Node, r *Result) {
		if node.Type == html.ElementNode && isWebsiteDetailBody(node) {
			isOk = true
			traverseDetailsBody(node.FirstChild, r)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseDetails(child, &result)
		}
	}

	traverseDetails(node, &result)

	return result, isOk
}

func parseListItem(node *html.Node) (Result, bool) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isWebsiteDetail(child) {
			return parseListItemDetails(child)
		}
	}

	return newEmptyResult(), false
}

func ParseResults(body string) ([]Result, error) {
	reader := strings.NewReader(body)
	doc, err := html.Parse(reader)

	if err != nil {
		return []Result{}, err
	}

	list, isListFound := findList(doc.FirstChild)

	if !isListFound {
		return []Result{}, &DomNodeNotFoundError{}
	}

	listItems := getListItems(list)

	if len(listItems) == 0 {
		return []Result{}, &DomNodeNotFoundError{}
	}

	results := []Result{}

	for _, listItemNode := range listItems {
		result, isOk := parseListItem(listItemNode)

		if isOk {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return results, &DomNodeNotFoundError{}
	}

	return results, nil
}
