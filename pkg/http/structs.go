package http

import (
    "wait4it/pkg/model"
)

type HttpCheck struct {
    URL            string // Url field renamed to URL
    Status         int
    Text           string
    FollowRedirect bool
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
    check := &HttpCheck{}
    check.BuildContext(*c)
    if err := check.ValidateURL(); err != nil { // Refactored Validate method to be consistent
        return nil, err
    }

    return check, nil
}

// ValidateURL ensures the URL field's consistency and correctness
func (h *HttpCheck) ValidateURL() error {
    // Assuming logic for URL validation is implemented here
    // This is a placeholder for actual validation logic
    return nil
}
