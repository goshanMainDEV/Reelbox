package httpserver

import "html/template"

type VideoMeta struct {
	Title    string
	ThumbURL string
	VideoURL string
}

var ogTemplate = template.Must(template.New("og").Parse(`<!DOCTYPE html>
<html prefix="og: https://ogp.me/ns#">
<head>
	<meta charset="utf-8">
	<meta property="og:type" content="video.other">
	<meta property="og:title" content="{{.Title}}">
	<meta property="og:image" content="{{.ThumbURL}}">
	<meta property="og:video" content="{{.VideoURL}}">
	<meta property="og:video:type" content="video/mp4">
</head>
<body>
	<video controls poster="{{.ThumbURL}}">
		<source src="{{.VideoURL}}" type="video/mp4">
	</video>
</body>
</html>`))