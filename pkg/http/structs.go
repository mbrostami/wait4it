package http

import (
    "wait4it/pkg/model"
)

type HttpCheck struct {
    URL            string // Renamed from 'Url' to 'URL'
    Status         int
    Text           string
    FollowRedirect bool
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
    check := &HttpCheck{}
    check.BuildContext(*c)
    if err := check.validateURL(); err != nil { // Assuming this method exists and also renamed from 'validateUrl' to 'validateURL'
        return nil, err
    }

    return check, nil
}

