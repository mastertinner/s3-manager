package s3manager_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mastertinner/s3manager/internal/app/s3manager"
	"github.com/mastertinner/s3manager/internal/app/s3manager/mocks"
	"github.com/matryer/is"
)

func TestHandleDeleteObject(t *testing.T) {
	t.Parallel()

	cases := []struct {
		it                   string
		removeObjectFunc     func(string, string) error
		expectedStatusCode   int
		expectedBodyContains string
	}{
		{
			it: "deletes an existing object",
			removeObjectFunc: func(string, string) error {
				return nil
			},
			expectedStatusCode:   http.StatusNoContent,
			expectedBodyContains: "",
		},
		{
			it: "returns error if there is an S3 error",
			removeObjectFunc: func(string, string) error {
				return errS3
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedBodyContains: http.StatusText(http.StatusInternalServerError),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.it, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			s3 := &mocks.S3Mock{
				RemoveObjectFunc: tc.removeObjectFunc,
			}

			req, err := http.NewRequest(http.MethodDelete, "/api/buckets/bucketName/objects/objectName", nil)
			is.NoErr(err)

			rr := httptest.NewRecorder()
			handler := s3manager.HandleDeleteObject(s3)

			handler.ServeHTTP(rr, req)

			is.Equal(tc.expectedStatusCode, rr.Code)                             // status code
			is.True(strings.Contains(rr.Body.String(), tc.expectedBodyContains)) // body
		})
	}
}
