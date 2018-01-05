package handlers_test

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyrusn/goTestHelper"
)

var testStudentLogin = func(t *testing.T) {
	for _, student := range studentList {
		w := httptest.NewRecorder()
		formString := fmt.Sprintf(`{"Username":"%s", "Password":"%s"}`, student.Username, student.Password)
		postForm := strings.NewReader(formString)
		req := httptest.NewRequest("POST", "/auth/students/login/", postForm)
		r.ServeHTTP(w, req)

		resp := w.Result()

		body, err := ioutil.ReadAll(resp.Body)
		assert.OK(t, err)

		mapToken[student.Username] = string(body)
	}
}

var testTeacherLogin = func(t *testing.T) {
	for _, teacher := range teacherList {
		w := httptest.NewRecorder()
		formString := fmt.Sprintf(`{"Username":"%s", "Password":"%s"}`, teacher.Username, teacher.Password)
		postForm := strings.NewReader(formString)
		req := httptest.NewRequest("POST", "/auth/teachers/login/", postForm)
		r.ServeHTTP(w, req)

		resp := w.Result()

		body, err := ioutil.ReadAll(resp.Body)
		assert.OK(t, err)

		mapToken[teacher.Username] = string(body)
	}
}

var testRefresh = func(t *testing.T) {
	for _, s := range studentList {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/auth/refresh/", nil)
		addJWT2Header(s.Username, req)
		r.ServeHTTP(w, req)
		body := parseBody(w, t)
		mapToken[s.Username] = string(body)
		// TODO: write a better to test the refreshed mapToken
	}
}
