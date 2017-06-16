package ast

// Node is the basic element of the AST.
// Interfaces embed Node should have 'Node' name suffix.
type Node interface {
	// Accept accepts Visitor to visit itself.
	// The returned node should replace original node.
	// ok returns false to stop visiting.
	//
	// Implementation of this method should first call visitor.Enter,
	// assign the returned node to its method receiver, if skipChildren returns true,
	// children should be skipped. Otherwise, call its children in particular order that
	// later elements depends on former elements. Finally, return visitor.Leave.
	Accept(v Visitor) (node Node, ok bool)
	// Text returns the original text of the element.
	Text() string
	// SetText sets original text to the Node.
	SetText(text string)
}

// Statement is an interface for SQL execution.
// NOTE: all Statement implementations must be safe for
// concurrent using by multiple goroutines.
// If the Exec method requires any Execution domain local data,
// they must be held out of the implementing instance.
type Statement interface {
	// OriginText gets the origin SQL text.
	OriginText() string

	// Exec executes SQL and gets a Recordset.

	// IsPrepared returns whether this statement is prepared statement.
	IsPrepared() bool
}

// Visitor visits a Node.
type Visitor interface {
	// Enter is called before children nodes are visited.
	// The returned node must be the same type as the input node n.
	// skipChildren returns true means children nodes should be skipped,
	// this is useful when work is done in Enter and there is no need to visit children.
	Enter(n Node) (node Node, skipChildren bool)
	// Leave is called after children nodes have been visited.
	// The returned node's type can be different from the input node if it is a ExprNode,
	// Non-expression node must be the same type as the input node n.
	// ok returns false to stop visiting.
	Leave(n Node) (node Node, ok bool)
}
