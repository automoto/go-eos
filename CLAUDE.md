## Code Quality
- Early Returns: Use early returns to reduce nesting
- Idiomatic Go: Code must be `make lint`/`go fmt` clean and use golang best practices
- Error Handling: Proper error handling; avoid panics unless truly fatal
- Reduce Code Nesting Where Possible: To ensure code readability, try to reduce code nesting (Nesting Depth) unless its needed.
- Keep code simple and concise. Try not to do overly complex or clever code unless its needed.
- Avoid verbose comments, only add comments where extra context is really needed.

## Testing Conventions

### TDD Workflow
- Always write failing tests BEFORE implementation
- Use AAA pattern: Arrange-Act-Assert
- One assertion per test when possible
- Test names describe behavior: "should_return_empty_when_no_items"
- Use testify "assert" https://pkg.go.dev/github.com/stretchr/testify/assert for test cases
- Use table based tests where appropriate to keep our tests concise

### Test-First Rules
- When I ask for a feature, write tests first
- Tests should FAIL initially (no implementation exists)
- Only after tests are written, implement minimal code to pass