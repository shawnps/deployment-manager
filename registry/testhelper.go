/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registry

// TODO(jackgr): Mock github repository service to test package and template registry implementations.

import (
	"github.com/kubernetes/deployment-manager/common"
	"github.com/kubernetes/deployment-manager/util"

	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type TestURLAndError struct {
	URL string
	Err error
}

type testGithubRegistryProvider struct {
	shortURL  string
	responses map[Type]TestURLAndError
}

type testGithubRegistry struct {
	githubRegistry
	responses map[Type]TestURLAndError
}

func NewTestGithubRegistryProvider(shortURL string, responses map[Type]TestURLAndError) GithubRegistryProvider {
	return testGithubRegistryProvider{
		shortURL:  util.TrimURLScheme(shortURL),
		responses: responses,
	}
}

func (tgrp testGithubRegistryProvider) GetGithubRegistry(cr common.Registry) (GithubRegistry, error) {
	trimmed := util.TrimURLScheme(cr.URL)
	if strings.HasPrefix(trimmed, tgrp.shortURL) {
		ghr, err := newGithubRegistry(cr.Name, trimmed, cr.Format, nil)
		if err != nil {
			panic(fmt.Errorf("cannot create a github registry: %s", err))
		}

		return &testGithubRegistry{
			githubRegistry: *ghr,
			responses:      tgrp.responses,
		}, nil
	}

	panic(fmt.Errorf("unknown registry: %v", cr))
}

func (tgr testGithubRegistry) ListTypes(regex *regexp.Regexp) ([]Type, error) {
	panic(fmt.Errorf("ListTypes should not be called in the test"))
}

func (tgr testGithubRegistry) GetDownloadURLs(t Type) ([]*url.URL, error) {
	result := tgr.responses[t]
	URL, err := url.Parse(result.URL)
	if err != nil {
		panic(err)
	}

	return []*url.URL{URL}, result.Err
}
