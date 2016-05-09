package composition

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func Test_ContentFetcher_FetchDefinitionHash(t *testing.T) {
	a := assert.New(t)
	tests := []struct {
		fd1 *FetchDefinition
		fd2 *FetchDefinition
		eq  bool
	}{
		{
			NewFetchDefinition("/foo"),
			NewFetchDefinition("/foo"),
			true,
		},
		{
			NewFetchDefinition("/foo"),
			NewFetchDefinition("/bar"),
			false,
		},
		{
			&FetchDefinition{
				URL:            "/foo",
				Timeout:        time.Second,
				RequestHeaders: http.Header{"Some": {"header"}},
				Required:       false,
			},
			&FetchDefinition{
				URL:            "/foo",
				Timeout:        time.Second * 42,
				RequestHeaders: http.Header{"Some": {"header"}},
				Required:       true,
			},
			true,
		},
		{
			&FetchDefinition{
				URL:            "/foo",
				RequestHeaders: http.Header{"Some": {"header"}},
			},
			&FetchDefinition{
				URL:            "/foo",
				RequestHeaders: http.Header{"Some": {"other header"}},
			},
			false,
		},
	}

	for _, t := range tests {
		if t.eq {
			a.Equal(t.fd1.Hash(), t.fd2.Hash())
		} else {
			a.NotEqual(t.fd1.Hash(), t.fd2.Hash())
		}
	}
}

func Test_ContentFetcher_FetchingWithDependency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a := assert.New(t)

	loader := NewMockContentLoader(ctrl)
	barFd := getFetchDefinitionMock(ctrl, loader, "/bar", nil, time.Millisecond*2)
	fooFd := getFetchDefinitionMock(ctrl, loader, "/foo", []*FetchDefinition{barFd}, time.Millisecond*2)
	bazzFd := getFetchDefinitionMock(ctrl, loader, "/bazz", []*FetchDefinition{barFd}, time.Millisecond)

	fetcher := NewContentFetcher()
	fetcher.contentLoader = loader

	fetcher.AddFetchJob(fooFd)
	fetcher.AddFetchJob(bazzFd)

	results := fetcher.WaitForResults()

	a.Equal(3, len(results))

	a.Equal("/foo", results[0].Def.URL)
	a.Equal("/bazz", results[1].Def.URL)
	a.Equal("/bar", results[2].Def.URL)
}

func getFetchDefinitionMock(ctrl *gomock.Controller, loaderMock *MockContentLoader, url string, requiredContent []*FetchDefinition, loaderBlocking time.Duration) *FetchDefinition {
	fd := NewFetchDefinition(url)
	fd.Timeout = time.Second * 42

	content := NewMockContent(ctrl)
	content.EXPECT().
		RequiredContent().
		Return(requiredContent)

	loaderMock.EXPECT().
		Load(fd.URL, fd.Timeout).
		Do(
			func(url string, timeout time.Duration) {
				time.Sleep(loaderBlocking)
			}).
		Return(content, nil)

	return fd
}