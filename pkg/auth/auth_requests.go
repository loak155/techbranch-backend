package auth

import "regexp"

var AuthRequests = []AuthRequest{
	{Mehtod: "GET", URL: regexp.MustCompile(`/docs`), Auth: false},

	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/articles$`), Auth: false},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/articles$`), Auth: true},
	{Mehtod: "PUT", URL: regexp.MustCompile(`/v1/articles$`), Auth: true},
	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/articles/[0-9]*$`), Auth: false},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/articles/[0-9]*$`), Auth: true},

	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/oauth/google/callback$`), Auth: false},
	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/oauth/google/login$`), Auth: false},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/refresh-token$`), Auth: false},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/signin$`), Auth: false},
	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/signin/user$`), Auth: true},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/signout$`), Auth: true},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/signup$`), Auth: false},

	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/articles/[0-9]*/bookmarks$`), Auth: false},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/articles/[0-9]*/bookmarks$`), Auth: true},
	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/articles/[0-9]*/bookmarks/count$`), Auth: false},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/bookmarks$`), Auth: true},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/users/[0-9]*/articles/[0-9]*/bookmarks$`), Auth: true},
	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/users/[0-9]*/bookmarks$`), Auth: true},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/users/[0-9]*/bookmarks$`), Auth: true},

	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/articles/[0-9]*/comments$`), Auth: false},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/articles/[0-9]*/comments$`), Auth: true},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/comments$`), Auth: true},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/comments/[0-9]*$`), Auth: true},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/users/[0-9]*/articles/[0-9]*/comments$`), Auth: true},
	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/users/[0-9]*/comments$`), Auth: true},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/users/[0-9]*/comments$`), Auth: true},

	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/users$`), Auth: true},
	{Mehtod: "POST", URL: regexp.MustCompile(`/v1/users$`), Auth: true},
	{Mehtod: "PUT", URL: regexp.MustCompile(`/v1/users$`), Auth: true},
	{Mehtod: "GET", URL: regexp.MustCompile(`/v1/users/[0-9]*$`), Auth: true},
	{Mehtod: "DELETE", URL: regexp.MustCompile(`/v1/users/[0-9]*$`), Auth: true},
}

var AuthMethods = map[string]bool{
	"/proto.ArticleService/CreateArticle": true,
	"/proto.ArticleService/GetArticle":    false,
	"/proto.ArticleService/ListArticles":  false,
	"/proto.ArticleService/UpdateArticle": true,
	"/proto.ArticleService/DeleteArticle": true,

	"/proto.AuthService/Signup":              false,
	"/proto.AuthService/Signin":              false,
	"/proto.AuthService/Signout":             false,
	"/proto.AuthService/RefreshToken":        true,
	"/proto.AuthService/GetSigninUser":       true,
	"/proto.AuthService/GetGoogleLoginURL":   false,
	"/proto.AuthService/GoogleLoginCallback": false,

	"/proto.BookmarkService/CreateBookmark":                     true,
	"/proto.BookmarkService/GetBookmarkCountByArticleID":        false,
	"/proto.BookmarkService/ListBookmarksByUserID":              true,
	"/proto.BookmarkService/ListBookmarksByArticleID":           false,
	"/proto.BookmarkService/DeleteBookmarkByUserIDAndArticleID": true,
	"/proto.BookmarkService/DeleteBookmarkByUserID":             true,
	"/proto.BookmarkService/DeleteBookmarkByArticleID":          true,

	"/proto.CommentService/CreateComment":                     true,
	"/proto.CommentService/ListCommentsByUserID":              true,
	"/proto.CommentService/ListCommentsByArticleID":           false,
	"/proto.CommentService/DeleteComment":                     true,
	"/proto.CommentService/DeleteCommentByUserIDAndArticleID": true,
	"/proto.CommentService/DeleteCommentByUserID":             true,
	"/proto.CommentService/DeleteCommentByArticleID":          true,

	"/proto.UserService/CreateUser": true,
	"/proto.UserService/GetUser":    true,
	"/proto.UserService/ListUsers":  true,
	"/proto.UserService/UpdateUser": true,
	"/proto.UserService/DeleteUser": true,
}
