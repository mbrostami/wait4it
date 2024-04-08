package http

import (
    "wait4it/pkg/model"
)

type HttpCheck struct {
    URL            string
    Status         int
    Text           string
    FollowRedirect bool
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
    check := &HttpCheck{}
    check.BuildContext(*c)
    if err := check.Validate(); err != nil { // This ensures existing methods like Validate are correctly refactored and consistent.
        return nil, err
    }

    return check, nil
}
