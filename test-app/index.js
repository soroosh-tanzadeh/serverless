import("./second")
response(
    '<h1>Click to show hello world</h1>' +
    '<button onclick="alert(\'Hello World\')">Click Me</button>',
    200,
    {
        "content-type": "text/html",
        "x-cache": true,
        "x-server": 1025
    }
)
