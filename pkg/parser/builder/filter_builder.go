package filter_builder

import "github.com/AbM7797/prisma-query-parser/pkg/parser"

// FilterBuilder helps parse nested filters safely
type FilterBuilder struct {
	current interface{}
	root    parser.Filter
	cursor  parser.Filter
	stack   []parser.Filter
}

// NewFilterBuilder creates a new builder instance
func NewFilterBuilder(f parser.Filter) *FilterBuilder {
	//root := make(utils.Filter)
	if f == nil {
		f = make(parser.Filter)
	}
	return &FilterBuilder{
		root:    f,
		cursor:  f,
		stack:   []parser.Filter{},
		current: f,
	}
}

// GetMap retrieves or initializes a nested map[string]interface{} by key and sets it as the current scope
func (b *FilterBuilder) GetMap(key string) *FilterBuilder {
	if b.current == nil {
		return b
	}

	// Assert current value as map[string]interface{}
	m, ok := b.current.(parser.Filter)
	if !ok {
		b.current = nil
		return b
	}

	// Get or initialize the nested value for the key
	val, exists := m[key]
	if !exists {
		nested := make(parser.Filter) // Initialize if not present
		m[key] = nested
		b.current = nested
		return b
	}

	// Ensure the value is a map[string]interface{}
	nested, ok := val.(parser.Filter)
	if !ok {
		nested = make(parser.Filter)
		m[key] = nested
	}

	b.current = nested
	return b
}

// GetString retrieves a string value by key
func (b *FilterBuilder) GetString(key string) *string {
	if b.current == nil {
		return nil
	}

	// Assert current value as map[string]interface{}
	m, ok := b.current.(parser.Filter)
	if !ok {
		return nil
	}

	val, exists := m[key]
	if !exists {
		return nil
	}

	// Assert the value as string
	strVal, ok := val.(string)
	if !ok {
		return nil
	}

	return &strVal
}

// Get retrieves the final value as an interface{}
func (b *FilterBuilder) Get() interface{} {
	return b.current
}

// GetFilter retrieves the current filter as utils.Filter (map[string]interface{})
func (b *FilterBuilder) GetFilter() parser.Filter {
	if b.current == nil {
		return nil
	}

	// Assert current value as map[string]interface{}
	m, ok := b.current.(parser.Filter)
	if !ok {
		return nil
	}
	// Explicitly cast to utils.Filter
	return parser.Filter(m)
}

// Reset resets the builder with a new filter
func (b *FilterBuilder) Reset(f parser.Filter) *FilterBuilder {
	b.current = f
	b.cursor = f
	b.root = f
	b.stack = []parser.Filter{}
	return b
}

// Add sets or merges a key-value pair into the current scope
func (b *FilterBuilder) AddExisting(key string, value interface{}) *FilterBuilder {
	if b.current == nil {
		return b
	}

	// Assert current value as map[string]interface{}
	m, ok := b.current.(parser.Filter)
	if !ok {
		return b
	}

	// Check if the key already exists
	if existing, exists := m[key]; exists {
		// If the existing value is a slice, append the new value
		if existingSlice, ok := existing.([]interface{}); ok {
			m[key] = append(existingSlice, value)
		} else {
			// If the existing value is not a slice, convert it to a slice
			m[key] = []interface{}{existing, value}
		}
	} else {
		// If the key doesn't exist, set it directly
		m[key] = value
	}
	return b
}

// Add sets a key-value pair at the current level
func (b *FilterBuilder) Add(key string, value interface{}) *FilterBuilder {
	b.cursor[key] = value
	return b
}

// Nest creates a new nested map at the given key and moves the cursor into it
func (b *FilterBuilder) Nest(key string) *FilterBuilder {
	nested := make(parser.Filter)
	b.cursor[key] = nested
	b.stack = append(b.stack, b.cursor) // Save current level
	b.cursor = nested                   // Move cursor to new level
	return b
}

// Up moves the cursor one level up in the nested structure
func (b *FilterBuilder) Up() *FilterBuilder {
	if len(b.stack) > 0 {
		// Pop the last element from the stack
		b.cursor = b.stack[len(b.stack)-1]
		b.stack = b.stack[:len(b.stack)-1]
	}
	return b
}

// Build returns the final filter structure
func (b *FilterBuilder) Build() parser.Filter {
	// Reset cursor to root
	b.cursor = b.root
	b.stack = []parser.Filter{} // Clear the stack
	return b.root
}
