package parser

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
	"time"
)

type DomNodeNotFoundError struct{}

func (*DomNodeNotFoundError) Error() string {
	return "Required HTML structure not present"
}

type Result struct {
	Url        string
	Title      string
	Summary    string
	AuthorName string
	UpdatedAt  int64
}

func (this *Result) HasAuthorName() bool {
	return this.AuthorName != ""
}

func (this *Result) HasUpdatedAt() bool {
	return this.UpdatedAt != 0
}

func NewEmptyResult() Result {
	return Result{
		Url:        "",
		Title:      "",
		Summary:    "",
		AuthorName: "",
		UpdatedAt:  0,
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
	return node.Data == "ol" && hasClass(&node.Attr, "websites")
}

func isWebsiteListItem(node *html.Node) bool {
	return node.Data == "li" && hasClass(&node.Attr, "websites__item")
}

func isWebsiteDetail(node *html.Node) bool {
	return node.Data == "details" && hasClass(&node.Attr, "website__details")
}

func isWebsiteDetailBody(node *html.Node) bool {
	return node.Data == "div" && hasClass(&node.Attr, "website__details__body")
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
	title := fmt.Sprint(node.FirstChild.Data)
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			url = attr.Val
		}
	}

	return url, title
}

func extractAuthorName(text string) string {
	parts := strings.Split(text, ",")

	if len(parts) == 1 {
		return ""
	}

	leading := strings.TrimSpace(parts[0])

	name, _ := strings.CutPrefix(leading, "By")

	return strings.TrimSpace(name)
}

func getAuthorName(node *html.Node) string {
	var authorName = ""
	var traverse func(node *html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "small" {
			authorName = extractAuthorName(node.FirstChild.Data)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(node)

	return authorName
}

func getUpdatedAt(node *html.Node) int64 {
	value := int64(0)

	var traverse func(node *html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "time" {
			for _, attr := range node.Attr {
				if attr.Key == "datetime" {
					dateTime, err := time.Parse("2006-01-02T15:04:05+00:00", attr.Val)

					if err == nil {
						value = dateTime.Unix()
					}
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(node)

	return value
}

func parseListItemDetails(node *html.Node) (Result, bool) {
	result := NewEmptyResult()
	isOk := false

	var traverseDetailsBody func(node *html.Node, r *Result)
	traverseDetailsBody = func(node *html.Node, r *Result) {
		if node.Type == html.ElementNode && node.Data == "a" {
			url, title := getLinkWithAnchor(node)

			r.Title = title
			r.Url = url
		}

		if node.Type == html.ElementNode && node.Data == "footer" {
			r.AuthorName = getAuthorName(node)
		}

		if node.Type == html.ElementNode && node.Data == "blockquote" {
			r.Summary = strings.TrimSpace(node.FirstChild.Data)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseDetailsBody(child, &result)
		}
	}

	var traverseDetails func(node *html.Node, r *Result)
	traverseDetails = func(node *html.Node, r *Result) {
		if node.Type == html.ElementNode && node.Data == "summary" {
			r.UpdatedAt = getUpdatedAt(node)
		}

		if node.Type == html.ElementNode && isWebsiteDetailBody(node) {
			isOk = true
			traverseDetailsBody(node, r)
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

	return NewEmptyResult(), false
}

func ParseResults(body string) ([]Result, error) {
	reader := strings.NewReader(body)
	doc, err := html.Parse(reader)

	if err != nil {
		return []Result{}, err
	}

	list, isListFound := findList(doc)

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
