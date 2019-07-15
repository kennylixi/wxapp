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
		url : "/mp/apply",
		data : $('#mpForm').serialize(),// 你的formid
		async : false,
		error : function(request) {
			parent.layer.alert("连接网络异常");
		},
		success : function(data) {
			if (data.code == 0) {
				parent.layer.msg("您的公众号已经被收录");
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