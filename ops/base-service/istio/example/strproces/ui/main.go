package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appConfig *Config
)

func page(c *gin.Context) {
	fmt.Println(c.Request.Header)
	c.Writer.Write([]byte(html))
}

func processer(c *gin.Context) {
	str := c.PostForm("str")

	if str == "" {
		c.String(http.StatusOK, "<h2>输入为空<h2>")
	}

	req := struct {
		Input string
	}{
		Input: str,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	httpReq, err := http.NewRequest("POST", appConfig.Service, bytes.NewReader(data))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	httpReq.Header = c.Request.Header
	resp, err := http.DefaultClient.Do(httpReq)
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}

	c.String(http.StatusOK, string(data))
}

func main() {
	engine := gin.Default()
	engine.GET("/ui", page)
	engine.POST("/ui", processer)
	http.ListenAndServe(":9090", engine)
}

const html = `<head>	
<script src='https://cdn.bootcss.com/jquery/3.3.1/jquery.min.js'></script>
</head>    
<html>
<body>
<h1>字符串处理</h1>
<input type="text" class="str" id="str">
<br>
<br>
<input id='ajax_btn' type='button' value='确认'>
<br><br>
<h3>处理结果</h3>
<div id='result'><h3></h3></div>
</body>
</html>
<script>
$(document).ready(function () { 
	 $('#ajax_btn').click(function () {
		var str =$("#str").val()
		var user=getQueryVariable("user")
		 $.ajax({
		   url: '/ui',
		   type: 'post',
		   dataType: 'html',
		   headers: {user:user},
		   data : {str: str},
		   success : function(data) {
			 $('#result').html(data);
		   },
		 });
	  });
});

function getQueryVariable(variable)
{
       var query = window.location.search.substring(1);
       var vars = query.split("&");
       for (var i=0;i<vars.length;i++) {
               var pair = vars[i].split("=");
               if(pair[0] == variable){return pair[1];}
       }
       return(false);
}
</script>`
