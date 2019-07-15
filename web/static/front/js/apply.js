$().ready(function() {
});

$('#Province').on('change', function() {
    selectLoad( $('#Province option:selected').val() )
});

function getCaptcha(){
    $.get("/captcha",function(data,status){
        if (status==="success"){
            $("#CaptchaId").val(data.CaptchaId);
            $("#cap").attr('src',data.Captcha);
        }
    });
}

function selectLoad(id) {
    var html = "";

    $('#City').empty();
    if(id===""){
        html += '<option value="">-请选择-</option>'
        $('#City').append(html);
        return
    }
    $.ajax({
        url: '/city/' + id,
        success: function (data) {
            $("#City").attr("style", "display:inline");
            //加载数据
            html += '<option value="">-请选择-</option>'
            for (var i = 0; i < data.length; i++) {
                html += '<option value="' + data[i].Id + '">' + data[i].Name + '</option>'
            }
            $('#City').append(html);
        }
    });
}

function save() {
	$.ajax({
		cache : true,
		type : "POST",
		url : "/apply",
		data : $('#signupForm').serialize(),// 你的formid
		async : false,
		error : function(request) {
			parent.layer.alert("连接网络异常");
		},
		success : function(data) {
			if (data.code == 0) {
				parent.layer.msg("您的小程序已经提交审核，请耐心等待");
                window.location.href="/"
            } else if (data.code == -1){
                var msg = "";
                for(var key in data.msg){
                    msg += data.msg[key];
                    msg += "<BR>";
                }
                parent.layer.msg(msg);
                getCaptcha();
            } else {
                parent.layer.msg(data.msg);
                getCaptcha();
            }

		}
	});

}

layui.use('upload', function(){
    var $ = layui.jquery,upload = layui.upload;

    //拖拽上传
    upload.render({
        elem: '#PicCoverElm',
        method:'POST'
        ,url: '/upload'
        ,done: function(data){
            if (data.code == 0) {
                $("#PicCoverPreview").attr("src", data.data.url);
                $("#PicCover").attr("value", data.data.url);
                $("#PicCoverElm").addClass("uploaded");
            } else {
                parent.layer.msg(data.msg);
            }
        }
    });

    upload.render({
        elem: '#PicQrcodeElm',
        method:'POST'
        ,url: '/upload'
        ,done: function(data){
            if (data.code == 0) {
                $("#PicQrcodePreview").attr("src", data.data.url);
                $("#PicQrcode").attr("value", data.data.url);
                $("#PicQrcodeElm").addClass("uploaded");
            } else {
                parent.layer.msg(data.msg);
            }
        }
    });

    upload.render({
        elem: '#Screenshot0Elm',
        method:'POST'
        ,url: '/upload'
        ,done: function(data){
            if (data.code == 0) {
                $(".screenshot-list").append('<div class="file uploaded screenshot-item"><img style="display:block" class="preview" src="'+data.data.url+'" width="100%" alt=""><input class="fileurl" type="hidden" name="Screenshot" value="'+data.data.url+'"><span onclick="deldiv(this)" class="del">删除</span></div>');
                if($(".screenshot-item").length>=5){
                    $(".add-screenshot").hide();
                }
            } else {
                parent.layer.msg(data.msg);
            }
        }
    });
});

function deldiv(e) {
    e.parentNode.remove();
    if($(".screenshot-item").length>=5){
        $(".add-screenshot").hide();
    } else {
        $(".add-screenshot").show();
    }
}

$("#cap").click(function () {
    getCaptcha()
});

var m=$("#header-nav-btn");
m.click(function(){
    var type=m.attr("data");
    if (type==1) {
        m.attr("data", 0);
        $("#header-nav").css("display","flex");
    } else {
        m.attr("data", 1);
        $("#header-nav").css("display","none");
    }
});