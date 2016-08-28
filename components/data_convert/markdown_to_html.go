package dataconvert

import (
	"strings"

	"github.com/russross/blackfriday"
)

const markdownToHTMLTemplate = `<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>Documents</title>
    <style media="screen">
        body {
            font-family: "Fira Sans", "Trebuchet MS", Arial, "Microsoft Yahei", sans-serif;
        }

        code {
            font-family: Consolas, "Fira Mono", "Liberation Mono", Courier, "Microsoft Yahei", monospace;
            color: red;
            padding: 2px 4px;
            font-size: 0.9em;
        }

        pre {
            border: 1px solid #CACACA;
            padding: 10px;
            overflow: auto;
            border-radius: 3px;
            background-color: #FAFAFB;
            color: #393939;
            margin: 2em 0px;
        }
        pre code {
            padding: 0;
        }

		table {
			border: 1px solid #CCC;
			border-collapse: collapse;
			width: 100%;
			margin-bottom: 20px;
		}

		th {
			border: 1px solid #CCC;
			background-color: aliceblue;
			white-space: nowrap;
			text-align:left;
			padding: 4px;
		}

		td {
			border: 1px solid #CCC;
			padding: 4px;
			word-break: break-word;
		}

        .body {
            max-width: 980px;
            margin: 0px auto;
        }
    </style>
</head>

<body>
    <div class="body">
		{{content}}
    </div>
</body>

</html>`

func markdownToHTML(input string) (string, error) {
	html := blackfriday.MarkdownCommon([]byte(input))

	return strings.Replace(markdownToHTMLTemplate, "{{content}}", string(html), 1), nil
}
