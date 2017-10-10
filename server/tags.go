package server

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
)

type qsNode struct {
	position *pongo2.Token
	key      pongo2.IEvaluator
	val      pongo2.IEvaluator
}

func (q *qsNode) Execute(ctx *pongo2.ExecutionContext, w pongo2.TemplateWriter) *pongo2.Error {
	r, _ := ctx.Public["request"].(*http.Request)
	if r == nil {
		return ctx.Error("request missing from context", q.position)
	}
	k, err := q.key.Evaluate(ctx)
	if err != nil {
		return err
	}
	v, err := q.val.Evaluate(ctx)
	if err != nil {
		return err
	}
	u := r.URL.Query()
	u.Set(k.String(), v.String())
	w.WriteString(fmt.Sprintf("%s?%s", r.URL.Path, u.Encode()))
	return nil
}

func qs(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	k, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	v, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	if arguments.Remaining() > 0 {
		return nil, arguments.Error("unexpected arguments", start)
	}
	return &qsNode{
		position: start,
		key:      k,
		val:      v,
	}, nil
}

func init() {
	pongo2.RegisterTag("qs", qs)
}
