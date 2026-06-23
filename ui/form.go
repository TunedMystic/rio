package ui

import "github.com/tunedmystic/rio/dom"

// Option is a single choice in a Select.
type Option struct {
	Value string
	Label string
}

// inputClasses is the shared field styling for text-like inputs.
const inputClasses = "block w-full rounded-[var(--radius-base)] border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-2 text-[var(--color-text)] shadow-sm transition placeholder:text-[var(--color-text-muted)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-primary)] focus:border-[var(--color-primary)]"

// Label renders a form label bound to a field id.
func Label(forID, text string, attrs ...dom.Node) dom.Node {
	children := make([]dom.Node, 0, len(attrs)+3)
	children = append(children, dom.Class("block text-[length:var(--font-size-sm)] font-medium text-[var(--color-text)] mb-1"), dom.For(forID))
	children = append(children, attrs...)
	children = append(children, dom.Text(text))
	return dom.Label(children...)
}

// FieldError renders small error text under a field, or nothing when msg is empty.
func FieldError(msg string) dom.Node {
	if msg == "" {
		return dom.Raw("")
	}
	return dom.P(
		dom.Class("text-[length:var(--font-size-sm)] text-[var(--color-danger)] mt-1"),
		dom.Text(msg),
	)
}

// TextField is a label + text input + optional error. The form workhorse.
func TextField(name, label, value, errMsg string, attrs ...dom.Node) dom.Node {
	input := make([]dom.Node, 0, len(attrs)+5)
	input = append(input, dom.Class(inputClasses), dom.Type("text"), dom.Id(name), dom.Name(name), dom.Value(value))
	input = append(input, attrs...)
	return dom.Div(
		dom.Class("mb-4"),
		Label(name, label),
		dom.Input(input...),
		FieldError(errMsg),
	)
}

// Textarea is a label + multi-line input + optional error.
func Textarea(name, label, value, errMsg string, attrs ...dom.Node) dom.Node {
	ta := make([]dom.Node, 0, len(attrs)+4)
	ta = append(ta, dom.Class(inputClasses), dom.Id(name), dom.Name(name))
	ta = append(ta, attrs...)
	ta = append(ta, dom.Text(value))
	return dom.Div(
		dom.Class("mb-4"),
		Label(name, label),
		dom.Textarea(ta...),
		FieldError(errMsg),
	)
}

// Select is a label + dropdown built from options, marking the selected value.
func Select(name, label string, options []Option, selected, errMsg string) dom.Node {
	return dom.Div(
		dom.Class("mb-4"),
		Label(name, label),
		dom.Select(
			dom.Class(inputClasses),
			dom.Id(name),
			dom.Name(name),
			dom.Map(options, func(o Option) dom.Node {
				opt := make([]dom.Node, 0, 3)
				opt = append(opt, dom.Value(o.Value))
				if o.Value == selected {
					opt = append(opt, dom.Selected())
				}
				opt = append(opt, dom.Text(o.Label))
				return dom.Option(opt...)
			}),
		),
		FieldError(errMsg),
	)
}

// Checkbox is a box + label.
func Checkbox(name, label string, checked bool, attrs ...dom.Node) dom.Node {
	return checkable("checkbox", name, label, "", checked, attrs...)
}

// Radio is a circle + label.
func Radio(name, label, value string, checked bool, attrs ...dom.Node) dom.Node {
	return checkable("radio", name, label, value, checked, attrs...)
}

func checkable(kind, name, label, value string, checked bool, attrs ...dom.Node) dom.Node {
	input := make([]dom.Node, 0, len(attrs)+5)
	input = append(input, dom.Class("h-4 w-4 rounded border-[var(--color-border)] text-[var(--color-primary)] accent-[var(--color-primary)] focus-visible:ring-2 focus-visible:ring-[var(--color-primary)] focus-visible:ring-offset-2"), dom.Type(kind), dom.Name(name))
	if value != "" {
		input = append(input, dom.Value(value))
	}
	if checked {
		input = append(input, dom.Checked())
	}
	input = append(input, attrs...)
	return dom.Label(
		dom.Class("inline-flex items-center gap-2 text-[var(--color-text)]"),
		dom.Input(input...),
		dom.Text(label),
	)
}
